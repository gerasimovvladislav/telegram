package telegram

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type B struct {
	client *WrappedClient

	userStorage UserStorage

	inlineCommandProcessorMap map[InlineCommandId]InlineCommandProcessor
	slashCommandProcessorMap  map[SlashCommandId]SlashCommandProcessor
	textProcessors            []TextProcessor
	textCommandConverters     []TextCommandConverter

	homeMarkup tgbotapi.InlineKeyboardMarkup
}

func NewBot(
	client *WrappedClient,
	userStorage UserStorage,
	inlineCommandProcessorMap map[InlineCommandId]InlineCommandProcessor,
	slashCommandProcessorMap map[SlashCommandId]SlashCommandProcessor,
	textProcessors []TextProcessor,
	textCommandConverters []TextCommandConverter,
) *B {
	return &B{
		client: WrapClient(client),

		userStorage: userStorage,

		inlineCommandProcessorMap: inlineCommandProcessorMap,
		slashCommandProcessorMap:  slashCommandProcessorMap,
		textProcessors:            textProcessors,
		textCommandConverters:     textCommandConverters,
	}
}

// Start запускает работу бота
func (b *B) Start(ctx context.Context) error {
	conf := tgbotapi.NewUpdate(0)
	conf.Timeout = 600

	slog.Info("Telegram Bot started")

	for {
		select {
		case <-ctx.Done():
			slog.Info("context is done, bot stop working")
			return nil
		default:
			var updates []tgbotapi.Update
			var err error

			doneCh := make(chan struct{})
			go func() {
				defer close(doneCh)
				for i := 0; i < 3; i++ {
					updates, err = b.client.GetUpdates(conf)
					if err != nil {
						seconds := 3
						slog.Error("can't get updates from api, will try next time after N seconds",
							err,
							slog.Int("seconds", seconds),
							slog.Int("attempt", i))
						time.Sleep(time.Second * time.Duration(seconds))
						continue
					}

					break
				}
			}()

			select {
			case <-ctx.Done():
				slog.Info("context is done, bot stop working on get updates step")
				return nil
			case <-doneCh:
			}

			if err != nil {
				return fmt.Errorf("can't get updates from api, after N attemts, last error: %w", err)
			}

			for _, update := range updates {
				if update.UpdateID >= conf.Offset {
					conf.Offset = update.UpdateID + 1
					go func(update tgbotapi.Update) {
						wrappedUpdate := WrapUpdate(&update)
						err = b.handleUpdate(ctx, wrappedUpdate)
						if err != nil {
							slog.Error("can't handle update", err)
						}
					}(update)
				}
			}
		}
	}
}

// handleUpdate обрабатывает обновление, которое пришло от апи-телеграмма
func (b *B) handleUpdate(ctx context.Context, update *Update) error {
	userId := update.UserId()
	if userId == 0 {
		slog.Warn("can't receive user id from bot update")

		return fmt.Errorf("can't receive user id from bot update")
	}

	logger := slog.With(slog.Int64("user_id", int64(userId)))
	logger.Debug("received update")
	if !update.Processable() {
		return fmt.Errorf("update is not processable")
	}

	defer func() {
		if msg := recover(); msg != nil {
			err := fmt.Errorf("%s", msg)
			slog.Error("router: panicked on handle update", err)
		}
	}()

	user, err := b.userStorage.FindById(userId)
	if err != nil {
		logger.Error("can't find user in storage", err)

		return err
	}
	if user == nil {
		user = NewAnonUser(userId)
		err = b.userStorage.Update(user)
		if err != nil {
			logger.Error("can't insert bot user in storage", err)

			return err
		}
	}

	var processorResult ProcessorResult
	var unresolvedCommandError *UnresolvedCommandError

	switch update.Type() {
	case MessageTypeInlineCommand:
		commandState := &InlineCommandState{}
		err = commandState.Decode(update.Raw().CallbackQuery.Data)
		if err != nil {
			logger.Error("can't unmarshal callback data from message: %w",
				err,
				slog.Int64("user_id", int64(userId)),
				slog.String("data", update.Raw().CallbackQuery.Data))
			return fmt.Errorf("can't unmarshal callback data from message: %w", err)
		}

		processorResult, err = b.tryExecuteCommand(commandState, user)
		if err != nil {
			if errors.As(err, &unresolvedCommandError) {
				logger.Info("can't execute command, wrong place",
					err,
					slog.Int("command_id", int(commandState.Id)))
			} else {
				logger.Error("can't execute command",
					err,
					slog.Int("command_id", int(commandState.Id)))
			}
			b.handleError(err, user)

			return err
		}
	case MessageTypeSlashCommand:
		var commandState *SlashCommandState
		commandState, err = NewSlashCommandState(update)
		if err != nil {
			logger.Error("can't create slash command state from message")
			b.handleError(err, user)

			return err
		}

		commandProcessor, ok := b.slashCommandProcessorMap[commandState.Id()]
		if !ok {
			msg := "can't find slash command processor for command"
			logger.Debug(msg)

			return fmt.Errorf(msg)
		}

		if len(commandProcessor.PlacesIn()) > 0 {
			isAvail := false
			for _, placeId := range commandProcessor.PlacesIn() {
				if placeId == user.Place.Id {
					isAvail = true
				}
			}

			if !isAvail {
				logger.Debug("user can't execute this command",
					slog.String("place_id", string(user.Place.Id)))

				return fmt.Errorf("user can't execute this command")
			}
		}

		processorResult, err = commandProcessor.Execute(commandState, user)
		if err != nil {
			if errors.As(err, &unresolvedCommandError) {
				logger.Info("can't execute inline command, wrong place", err)
			} else {
				logger.Error("can't execute inline command", err)
			}
			b.handleError(err, user)

			return err
		}
	case MessageTypeTextCommand:
		var textCommandProcessorToUse TextCommandConverter
		for _, cmd := range b.textCommandConverters {
			if cmd.CanConvert(update.Raw().Message.Text) {
				textCommandProcessorToUse = cmd
				break
			}
		}

		if textCommandProcessorToUse != nil {
			slog.Debug("simple message with text was a 'text command', converted into inline command",
				slog.Int("commandId", int(textCommandProcessorToUse.InlineCommandId())),
				slog.String("text", update.Raw().Message.Text))

			var result *TextCommandProcessorResult
			result, err = textCommandProcessorToUse.Process(update.Raw().Message.Text, user)
			if err != nil {
				slog.Error("can't convert text command", err)
				b.handleError(err, user)

				return err
			}

			for _, msg := range result.Msgs() {
				_, err = b.client.Send(msg)
				if err != nil {
					slog.Error("can't send messages", err)
				}
			}

			processorResult, err = b.tryExecuteCommand(result.CommandCall(), user)
			if err != nil {
				if errors.As(err, &unresolvedCommandError) {
					slog.Info("can't execute command after translated from text command, wrong place", err)
				} else {
					slog.Error("can't execute command after translated from text command", err)
				}
				b.handleError(err, user)

				return err
			}
		}
	case MessageTypeMessage:
		var processorToUse TextProcessor
		for _, p := range b.textProcessors {
			if len(p.PlacesIn()) == 0 {
				processorToUse = p
				break
			}

			for _, st := range p.PlacesIn() {
				if st == user.Place.Id {
					processorToUse = p
					break
				}
			}

			if processorToUse != nil {
				break
			}
		}

		if processorToUse != nil {
			processorResult, err = processorToUse.Process(update, user)
			if err != nil {
				b.handleError(err, user)
				if errors.As(err, &unresolvedCommandError) {
					logger.Info("can't process message, the command was called from the wrong place", err)
				} else {
					logger.Error("can't process message", err)
				}

				return err
			}
		} else {
			logger.Debug("processor not found")
			b.handleError(nil, user)
		}
	}

	if processorResult == nil {
		//TODO: подумать надо ли писать лог
		//return fmt.Errorf("nothing to process, processorResult == nil")
		return nil
	}

	err = b.handleProcessorResult(ctx, processorResult, user)
	if err != nil {
		if errors.As(err, &unresolvedCommandError) {
			slog.Info("can't handle processor result, the command was called from the wrong place", err)
		} else {
			slog.Error("can't handle processor result", err)
		}
		return err
	}

	return nil
}

func (b *B) handleProcessorResult(
	ctx context.Context,
	result ProcessorResult,
	user *User,
) error {
	for _, msg := range result.Msgs() {
		_, err := b.client.Send(msg)
		if err != nil {
			return fmt.Errorf("can't send msg: %w", err)
		}
	}

	//TODO: передумать
	u, err := b.userStorage.FindById(user.ID)
	if err != nil {
		return fmt.Errorf("can't get bot user")
	}

	if u != nil {
		if result.NewPlace() != nil && !result.NewPlace().Eq(u.Place) {
			// set new state
			slog.Debug("set new place for user",
				slog.String("old_place_id", string(u.Place.Id)),
				slog.String("new_place_id", string(result.NewPlace().Id)))

			u.Place = result.NewPlace()
			err = b.userStorage.Update(u)
			if err != nil {
				return fmt.Errorf("can't update user state: %w", err)
			}
		}

		if result.CommandCall() != nil {
			var processorResult *InlineCommandProcessorResult
			processorResult, err = b.tryExecuteCommand(result.CommandCall(), u)
			if err != nil {
				return fmt.Errorf("can't call next process command: %w", err)
			}

			return b.handleProcessorResult(ctx, processorResult, u)
		}
	} else {
		slog.Debug("bot user was deleted in previous actions, can't do anything additional")
	}

	return nil
}

func (b *B) handleError(processorErr error, user *User) {
	var msg tgbotapi.MessageConfig
	var unresolvedErr *UnresolvedCommandError

	errorText := "Произошла ошибка, повторите запрос позднее. Доступные команды можно получить в разделе [/help](tg:///help)"
	if errors.As(processorErr, &unresolvedErr) {
		errorText = "Завершите предыдущее действие"
	}

	if user.Place.Id == PlaceIdEmpty {
		msg = tgbotapi.NewMessage(int64(user.ID), errorText)
		msg.ParseMode = tgbotapi.ModeMarkdown
		msg.ReplyMarkup = b.homeMarkup
	} else {
		msg = tgbotapi.NewMessage(int64(user.ID), errorText)
		msg.ParseMode = tgbotapi.ModeMarkdown
		msg.ReplyMarkup = b.homeMarkup
	}

	_, err := b.client.Send(msg)
	if err != nil {
		slog.Error("can't send message", processorErr)
	}
}

func (b *B) tryExecuteCommand(
	commandState *InlineCommandState,
	user *User,
) (*InlineCommandProcessorResult, error) {
	commandProcessor, ok := b.inlineCommandProcessorMap[commandState.Id]
	if !ok {
		return nil, fmt.Errorf("can't find command processor for command '%d'", commandState.Id)
	}

	if len(commandProcessor.PlacesIn()) > 0 {
		placeAllowed := false
		for _, placeId := range commandProcessor.PlacesIn() {
			if placeId == user.Place.Id {
				placeAllowed = true
				break
			}
		}

		if !placeAllowed {
			return nil, NewUnresolvedCommandError(
				fmt.Errorf("command %d:'%s' can't be executed on state '%s'",
					commandProcessor.Id(),
					"//TODO", //TODO: выводить имя команды
					user.Place.Id),
			)
		}
	}

	processorResult, err := commandProcessor.Execute(commandState, user)
	if err != nil {
		b.handleError(err, user)
		return nil, fmt.Errorf("can't execute command: %w", err)
	}

	return processorResult, nil
}

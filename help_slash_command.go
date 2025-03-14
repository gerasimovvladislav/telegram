package telegram

import (
	"gopkg.in/telegram-bot-api.v4"
)

func newHelpSlashCommandProcessor() *HelpSlashCommandProcessor {
	return &HelpSlashCommandProcessor{}
}

// HelpSlashCommandProcessor обработчик команды /help
type HelpSlashCommandProcessor struct {
}

func (p *HelpSlashCommandProcessor) Id() SlashCommandId {
	return SlashCommandIdHelp
}

func (p *HelpSlashCommandProcessor) PlacesIn() []PlaceId {
	return nil
}

func (p *HelpSlashCommandProcessor) PlaceOut() PlaceId {
	return ""
}

// Execute получить инфо по командам бота
func (p *HelpSlashCommandProcessor) Execute(_ *SlashCommandState, user *User) (*SlashCommandProcessorResult, error) {
	msgBody := `
		Доступные команды:

		[/start](tg:///start) - начать диалог с ботом

		[/help](tg:///help) - список доступных команд
		
		[/myid](tg:///myid) - вывести информацию о пользователе

		[/reset](tg:///reset) - произвести сброс данных пользователя
	`

	msg := tgbotapi.NewMessage(int64(user.ID), msgBody)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)

	return NewSlashCommandProcessorResult([]tgbotapi.Chattable{msg}, nil), nil
}

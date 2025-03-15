package telegram

import (
	"errors"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SlashCommandId string

const (
	// SlashCommandIdStart начало работы с ботом /start
	SlashCommandIdStart SlashCommandId = "start"
	// SlashCommandIdMyId информация о пользователе /myid
	SlashCommandIdMyId SlashCommandId = "myid"
	// SlashCommandIdReset сброс настроек пользователя /reset
	SlashCommandIdReset SlashCommandId = "reset"
	// SlashCommandIdHelp раздел помощи /help
	SlashCommandIdHelp SlashCommandId = "help"
)

func NewSlashCommandState(update *Update) (*SlashCommandState, error) {
	if update.Raw().Message == nil || !update.Raw().Message.IsCommand() {
		return nil, errors.New("update is not valid slash command")
	}
	return &SlashCommandState{update: update}, nil
}

type SlashCommandState struct {
	update *Update
}

func (c *SlashCommandState) Id() SlashCommandId {
	return SlashCommandId(c.update.Raw().Message.Command())
}

func (c *SlashCommandState) Args() []string {
	return strings.SplitAfterN(strings.TrimLeft(c.update.Raw().Message.Text, "/"+c.update.Raw().Message.Command()+" "), " ", 1)
}

type SlashCommandProcessor interface {
	Id() SlashCommandId
	// PlacesIn возвращает разрешенные для обработчика плейсы, откуда можно вызвать обработчик
	PlacesIn() []PlaceId
	// PlaceOut PlacesOut возвращает разрешенные плейсы для следующих обработчиков
	PlaceOut() PlaceId
	Execute(state *SlashCommandState, user *User) (*SlashCommandProcessorResult, error)
}

func NewSlashCommandProcessorResult(
	msgs []tgbotapi.Chattable,
	newPlace *Place,
) *SlashCommandProcessorResult {
	return &SlashCommandProcessorResult{
		msgs:     msgs,
		newPlace: newPlace,
	}
}

type SlashCommandProcessorResult struct {
	// Список сообщений, которые надо отправить пользователю
	msgs []tgbotapi.Chattable
	// Новое состояние, которое необходимо проставить пользователю
	newPlace *Place
}

func (r *SlashCommandProcessorResult) Msgs() []tgbotapi.Chattable {
	return r.msgs
}

func (r *SlashCommandProcessorResult) NewPlace() *Place {
	return r.newPlace
}

func (r *SlashCommandProcessorResult) CommandCall() *InlineCommandState {
	return nil // для команд не предусмотрен возврат вызова других команд, пока...
}

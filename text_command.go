package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TextCommandId string

type TextCommandConverter interface {
	Id() TextCommandId
	InlineCommandId() InlineCommandId
	CanConvert(text string) bool
	Process(text string, user *User) (*TextCommandProcessorResult, error)
}

func NewTextCommandProcessorResult(
	msgs []tgbotapi.Chattable,
	commandCall *InlineCommandState,
) *TextCommandProcessorResult {
	return &TextCommandProcessorResult{
		msgs:        msgs,
		commandCall: commandCall,
	}
}

type TextCommandProcessorResult struct {
	// Список сообщений, которые надо отправить пользователю
	msgs []tgbotapi.Chattable
	// Команда, которую надо вызвать
	commandCall *InlineCommandState
}

func (r *TextCommandProcessorResult) Msgs() []tgbotapi.Chattable {
	return r.msgs
}

func (r *TextCommandProcessorResult) NewPlace() *Place {
	return nil
}

func (r *TextCommandProcessorResult) CommandCall() *InlineCommandState {
	return r.commandCall
}

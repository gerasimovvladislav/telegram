package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TextProcessor interface {
	PlacesIn() []PlaceId
	PlacesOut() []PlaceId
	Process(update *Update, user *User) (*TextProcessorResult, error)
}

func NewTextProcessorResult(
	msgs []tgbotapi.Chattable,
	newPlace *Place,
	commandCall *InlineCommandState,
) *TextProcessorResult {
	return &TextProcessorResult{
		msgs:        msgs,
		newPlace:    newPlace,
		commandCall: commandCall,
	}
}

type TextProcessorResult struct {
	// Список сообщений, которые надо отправить пользователю
	msgs []tgbotapi.Chattable
	// Новое состояние, которое необходимо проставить пользователю
	newPlace *Place
	// Команда, которую необходимо выполнить. Данная опция нужна для случаев, когда после одного процесс нужно сразу
	// запускать следующий процесс. Например, после авторизации клиента, надо сразу предложить ему выбрать компании, а
	// не выводить му кнопку типа "выберите компанию"
	commandCall *InlineCommandState
}

func (r *TextProcessorResult) Msgs() []tgbotapi.Chattable {
	return r.msgs
}

func (r *TextProcessorResult) NewPlace() *Place {
	return r.newPlace
}

func (r *TextProcessorResult) CommandCall() *InlineCommandState {
	return r.commandCall
}

type ProcessorResult interface {
	// Msgs возвращает список сообщений, которые надо отправить пользователю
	Msgs() []tgbotapi.Chattable
	// NewPlace возвращает новое состояние, которое необходимо проставить пользователю
	NewPlace() *Place
	// CommandCall возвращает команду, которую необходимо выполнить. Данная опция нужна для случаев, когда после одного процесс нужно сразу
	// запускать следующий процесс. Например, после авторизации клиента, надо сразу предложить ему выбрать компании, а
	// не выводить му кнопку типа "выберите компанию"
	CommandCall() *InlineCommandState
}

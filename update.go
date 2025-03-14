package telegram

import (
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// Update Обрамление обновления бота телеграмма. Облегчает работу с этим обновлением
type Update struct {
	messageType MessageType
	update      *tgbotapi.Update
}

// Processable возможно ли обработать данное обновление
func (u *Update) Processable() bool {
	return u.IsInlineButtonPressed() || u.update.Message != nil
}

// UserId выделяет из обновления идентификатор пользователя
func (u *Update) UserId() UserId {
	if u.IsInlineButtonPressed() {
		return UserId(u.update.CallbackQuery.From.ID)
	}

	if u.update.EditedMessage != nil {
		return UserId(u.update.EditedMessage.From.ID)
	}

	if u.update.Message != nil {
		return UserId(u.update.Message.From.ID)
	}

	return 0
}

// IsMessageWithText удостоверяется, что обновление является сообщением с текстом
func (u *Update) IsMessageWithText() bool {
	if u.update.Message == nil {
		return false
	}

	if u.update.Message.Text == "" {
		return false
	}

	if u.update.Message.IsCommand() {
		return false
	}

	return true
}

// IsInlineButtonPressed проверяет, что обновление является нажатием кнопки
func (u *Update) IsInlineButtonPressed() bool {
	return u.update.CallbackQuery != nil
}

// Raw возвращает обновление, которые пришло из библиотеки телеграмма
func (u *Update) Raw() *tgbotapi.Update {
	return u.update
}

// Type возвращает тип сообщения
func (u *Update) Type() MessageType {
	return u.messageType
}

package telegram

import (
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// WrapUpdate оборачивает обновление библиотеки телеграмма, дает доступ к вспомогательным методам
func WrapUpdate(update *tgbotapi.Update) *Update {
	wrappedUpdate := &Update{update: update}

	wrappedUpdate.messageType = MessageTypeMessage
	switch {
	case wrappedUpdate.IsInlineButtonPressed():
		wrappedUpdate.messageType = MessageTypeInlineCommand
	case wrappedUpdate.Raw().Message != nil && wrappedUpdate.Raw().Message.IsCommand():
		wrappedUpdate.messageType = MessageTypeSlashCommand
	case wrappedUpdate.Raw().Message != nil && wrappedUpdate.Raw().Message.Text != "":
		wrappedUpdate.messageType = MessageTypeTextCommand
	}

	return wrappedUpdate
}

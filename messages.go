package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// NewMessage creates new message
func NewMessage(chatID int64, text string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown

	return msg
}

// ReplyMessage creates new reply message
func ReplyMessage(chatID int64, replyID int, text string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyToMessageID = replyID
	msg.ParseMode = tgbotapi.ModeHTML

	return msg
}

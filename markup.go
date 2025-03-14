package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func NewButton(text string, inlineCommandId InlineCommandId) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(
		text,
		NewSimpleInlineCommand(inlineCommandId).State().Encode(),
	)
}

func NewButtonWithMeta(text string, inlineCommandId InlineCommandId, meta []string) tgbotapi.InlineKeyboardButton {
	state := NewSimpleInlineCommand(inlineCommandId).State()
	state.Strings = meta

	return tgbotapi.NewInlineKeyboardButtonData(text, state.Encode())
}

func NewRow(buttons ...tgbotapi.InlineKeyboardButton) []tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardRow(buttons...)
}

func NewMarkup(rows ...[]tgbotapi.InlineKeyboardButton) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

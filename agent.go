package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewAgent(token string) (*tgbotapi.BotAPI, error) {
	tgBotAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("can't init telegram bot: %w", err)
	}
	tgBotAPI.Debug = false

	return tgBotAPI, nil
}

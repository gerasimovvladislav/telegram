package telegram

import (
	"fmt"
	"net/http"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func NewAgent(token string) (*tgbotapi.BotAPI, error) {
	tgBotAPI, err := tgbotapi.NewBotAPIWithClient(token, &http.Client{})
	if err != nil {
		return nil, fmt.Errorf("can't init telegram bot: %w", err)
	}
	tgBotAPI.Debug = false

	return tgBotAPI, nil
}

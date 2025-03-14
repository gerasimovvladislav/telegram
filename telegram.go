package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Agent interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	StopReceivingUpdates()
	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
	GetUpdates(config tgbotapi.UpdateConfig) ([]tgbotapi.Update, error)
}

// WrappedClient оборачивает клиент библиотеки телеграмма, дает доступ к вспомогательным методам
type WrappedClient struct {
	Agent
}

func WrapClient(agent Agent) *WrappedClient {
	return &WrappedClient{
		Agent: agent,
	}
}

// Send метод обертка
func (w *WrappedClient) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	message, err := w.Agent.Send(c)

	return message, err
}

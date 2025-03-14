package processors

import (
	tgbotapi "gopkg.in/telegram-bot-api.v4"

	"gitlab.com/vladislavgerasimov/telegram"
	"gitlab.com/vladislavgerasimov/telegram/examples/tgbot/internal/slash"
)

type Hello struct{}

func NewHello() *Hello {
	return &Hello{}
}

func (h *Hello) Id() telegram.SlashCommandId {
	return slash.CommandIdHello
}

func (h *Hello) PlacesIn() []telegram.PlaceId {
	return nil
}

func (h *Hello) PlaceOut() telegram.PlaceId {
	return ""
}

func (h *Hello) Execute(_ *telegram.SlashCommandState, user *telegram.User) (*telegram.SlashCommandProcessorResult, error) {
	msg := tgbotapi.NewMessage(
		int64(user.ID),
		"Hello, World!")
	msg.ParseMode = tgbotapi.ModeMarkdown

	return telegram.NewSlashCommandProcessorResult([]tgbotapi.Chattable{msg}, nil), nil
}

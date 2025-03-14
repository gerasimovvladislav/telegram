package telegram

import (
	"fmt"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func newMyIdSlashCommandProcessor() *MyIdSlashCommandProcessor {
	return &MyIdSlashCommandProcessor{}
}

type MyIdSlashCommandProcessor struct {
}

func (m *MyIdSlashCommandProcessor) Id() SlashCommandId {
	return SlashCommandIdMyId
}

func (m *MyIdSlashCommandProcessor) PlacesIn() []PlaceId {
	return nil
}

func (m *MyIdSlashCommandProcessor) PlaceOut() PlaceId {
	return ""
}

func (m *MyIdSlashCommandProcessor) Execute(_ *SlashCommandState, user *User) (*SlashCommandProcessorResult, error) {
	msg := tgbotapi.NewMessage(
		int64(user.ID),
		fmt.Sprintf("ID: %d;\nPlace: '%s';\n", user.ID, user.Place.Id),
	)
	msg.ParseMode = tgbotapi.ModeMarkdown

	return NewSlashCommandProcessorResult([]tgbotapi.Chattable{msg}, nil), nil
}

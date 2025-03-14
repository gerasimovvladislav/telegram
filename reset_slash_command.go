package telegram

import (
	"fmt"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func newResetSlashCommandProcessor() *ResetSlashCommandProcessor {
	return &ResetSlashCommandProcessor{}
}

type ResetSlashCommandProcessor struct {
}

func (m *ResetSlashCommandProcessor) Id() SlashCommandId {
	return SlashCommandIdReset
}

func (m *ResetSlashCommandProcessor) PlacesIn() []PlaceId {
	return nil
}

func (m *ResetSlashCommandProcessor) PlaceOut() PlaceId {
	return ""
}

func (m *ResetSlashCommandProcessor) Execute(_ *SlashCommandState, user *User) (*SlashCommandProcessorResult, error) {
	msg := tgbotapi.NewMessage(
		int64(user.ID),
		fmt.Sprintf("//TODO: Reset: ID: %d;\nPlace: '%s';\n", user.ID, user.Place.Id),
	)
	msg.ParseMode = tgbotapi.ModeMarkdown

	return NewSlashCommandProcessorResult([]tgbotapi.Chattable{msg}, nil), nil
}

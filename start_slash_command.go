package telegram

import (
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func newStartSlashCommandProcessor() *StartSlashCommandProcessor {
	return &StartSlashCommandProcessor{}
}

type StartSlashCommandProcessor struct {
}

func (p *StartSlashCommandProcessor) Id() SlashCommandId {
	return SlashCommandIdStart
}

func (p *StartSlashCommandProcessor) PlacesIn() []PlaceId {
	return nil
}

func (p *StartSlashCommandProcessor) PlaceOut() PlaceId {
	return ""
}

func (p *StartSlashCommandProcessor) Execute(_ *SlashCommandState, user *User) (*SlashCommandProcessorResult, error) {
	return NewSlashCommandProcessorResult([]tgbotapi.Chattable{p.newDefaultHelloMessage(user.ID)}, nil), nil
}

func (p *StartSlashCommandProcessor) newDefaultHelloMessage(userId UserId) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(
		int64(userId),
		"Привет, я *Бот* 👨‍💻!")
	msg.ParseMode = tgbotapi.ModeMarkdown

	return msg
}

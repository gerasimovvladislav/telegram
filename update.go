package telegram

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Update wraps tgbotapi.Update
type Update struct {
	Original *tgbotapi.Update `json:"raw,omitempty"`
}

// WrapUpdate returns wrapped update
func WrapUpdate(update *tgbotapi.Update) *Update {
	return &Update{Original: update}
}

// Processable checks if update can be processed
func (u *Update) Processable() bool {
	return u.IsPost() || u.IsInlineCommand() || u.IsSlashCommand() || u.IsTextMessage()
}

func (u *Update) IsPost() bool {
	return u.Original.ChannelPost != nil && u.Original.ChannelPost.ReplyToMessage != nil
}

// IsTextMessage checks if update is text message
func (u *Update) IsTextMessage() bool {
	if u.Original.Message == nil {
		return false
	}

	if u.Original.Message.Text == "" {
		return false
	}

	if u.Original.Message.IsCommand() {
		return false
	}

	return true
}

// IsInlineCommand checks if update is inline command
func (u *Update) IsInlineCommand() bool {
	return u.Original.CallbackQuery != nil
}

// IsSlashCommand c

func (u *Update) IsSlashCommand() bool {
	return u.Original.Message != nil && strings.HasPrefix(u.Original.Message.Text, "/")
}

// UserID returns user id
func (u *Update) UserID() UserID {
	if u.IsPost() {
		return UserID(u.Original.ChannelPost.From.ID)
	}

	if u.IsInlineCommand() {
		return UserID(u.Original.CallbackQuery.From.ID)
	}

	if u.Original.EditedMessage != nil {
		return UserID(u.Original.EditedMessage.From.ID)
	}

	if u.Original.Message != nil {
		return UserID(u.Original.Message.From.ID)
	}

	return 0
}

// Raw returns raw update
func (u *Update) Raw() *tgbotapi.Update {
	return u.Original
}

// Type returns update type
func (u *Update) Type() MessageType {
	updateType := UnknownUpdate

	switch {
	case u.IsInlineCommand():
		updateType = InlineCommand
	case u.IsSlashCommand():
		updateType = SlashCommand
	case u.IsTextMessage():
		updateType = TextMessage
	}

	return updateType
}

// SlashCommandID returns slash command id
func (u *Update) SlashCommandID() string {
	return u.Raw().Message.Command()
}

// SlashArgs returns slash command args
func (u *Update) SlashArgs() []string {
	return strings.SplitAfterN(
		strings.TrimLeft(
			u.Raw().Message.Text,
			"/"+u.Raw().Message.Command()+" ",
		),
		" ", -1)
}

// InlineCommandID returns inline command id
func (u *Update) InlineCommandID() string {
	parts := strings.Split(u.Raw().CallbackQuery.Data, "#")
	if len(parts) == 0 {
		return ""
	}

	return parts[0]
}

// InlineArgs returns inline command args
func (u *Update) InlineArgs() []string {
	parts := strings.Split(u.Raw().CallbackQuery.Data, "#")
	if len(parts) != 2 {
		return nil
	}

	args := strings.Split(parts[1], ";")
	if len(args) == 1 && args[0] == "" {
		return nil
	}

	return args
}

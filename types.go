package telegram

type UserID int64

func (uid UserID) Int64() int64 {
	return int64(uid)
}

type MessageType string

func (mt MessageType) String() string {
	return string(mt)
}

const (
	InlineCommand MessageType = "inline_command"
	SlashCommand  MessageType = "slash_command"
	TextMessage   MessageType = "text_message"
	UnknownUpdate MessageType = "unknown_update"
)

package telegram

type UserId int64

type MessageType string

func (mt MessageType) String() string {
	return string(mt)
}

const (
	MessageTypeInlineCommand MessageType = "inline_command"
	MessageTypeSlashCommand  MessageType = "slash_command"
	MessageTypeTextCommand   MessageType = "text_command"
	MessageTypeMessage       MessageType = "message"
)

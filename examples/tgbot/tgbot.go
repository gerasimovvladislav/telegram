package tgbot

import (
	"gitlab.com/vladislavgerasimov/telegram"
)

func Init(token string, IDs []int64,
	inlineCommandProcessorMap map[telegram.InlineCommandId]telegram.InlineCommandProcessor,
	slashCommandProcessorMap map[telegram.SlashCommandId]telegram.SlashCommandProcessor,
	textProcessors []telegram.TextProcessor,
	textCommandConverters []telegram.TextCommandConverter,
) (*telegram.B, error) {
	users := make([]*telegram.User, 0, len(IDs))
	for _, id := range IDs {
		users = append(users, telegram.NewAnonUser(telegram.UserId(id)))
	}
	storage := telegram.NewUsers(users)

	return telegram.Init(
		token,
		storage,
		inlineCommandProcessorMap,
		slashCommandProcessorMap,
		textProcessors,
		textCommandConverters,
	)
}

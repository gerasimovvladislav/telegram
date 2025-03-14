package telegram

import (
	"fmt"
)

func Init(
	token string,
	users UserStorage,
	inlineCommandProcessorMap map[InlineCommandId]InlineCommandProcessor,
	slashCommandProcessorMap map[SlashCommandId]SlashCommandProcessor,
	textProcessors []TextProcessor,
	textCommandConverters []TextCommandConverter,
) (*B, error) {
	agent, err := NewAgent(token)
	if err != nil {
		return nil, fmt.Errorf("error init agent: %w", err)
	}

	if start, ok := slashCommandProcessorMap[SlashCommandIdStart]; !ok || start == nil {
		slashCommandProcessorMap[SlashCommandIdStart] = newStartSlashCommandProcessor()
	}
	if help, ok := slashCommandProcessorMap[SlashCommandIdHelp]; !ok || help == nil {
		slashCommandProcessorMap[SlashCommandIdHelp] = newHelpSlashCommandProcessor()
	}

	if myid, ok := slashCommandProcessorMap[SlashCommandIdMyId]; !ok || myid == nil {
		slashCommandProcessorMap[SlashCommandIdMyId] = newMyIdSlashCommandProcessor()
	}

	if reset, ok := slashCommandProcessorMap[SlashCommandIdReset]; !ok || reset == nil {
		slashCommandProcessorMap[SlashCommandIdReset] = newResetSlashCommandProcessor()
	}

	if len(slashCommandProcessorMap) == 0 {
		slashCommandProcessorMap = map[SlashCommandId]SlashCommandProcessor{
			SlashCommandIdStart: newStartSlashCommandProcessor(),
			SlashCommandIdHelp:  newHelpSlashCommandProcessor(),
			SlashCommandIdMyId:  newMyIdSlashCommandProcessor(),
			SlashCommandIdReset: newResetSlashCommandProcessor(),
		}
	}

	return NewBot(
		WrapClient(agent),
		users,
		inlineCommandProcessorMap,
		slashCommandProcessorMap,
		textProcessors,
		textCommandConverters,
	), nil
}

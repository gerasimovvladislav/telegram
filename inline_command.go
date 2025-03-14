package telegram

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type SimpleInlineCommand struct {
	id InlineCommandId
}

func NewSimpleInlineCommand(id InlineCommandId) *SimpleInlineCommand {
	return &SimpleInlineCommand{
		id: id,
	}
}

func (c *SimpleInlineCommand) Id() InlineCommandId {
	return c.id
}

func (c *SimpleInlineCommand) State() *InlineCommandState {
	return &InlineCommandState{Id: c.id}
}

type InlineCommandProcessor interface {
	Id() InlineCommandId
	PlacesIn() []PlaceId
	PlaceOut() PlaceId
	Execute(state *InlineCommandState, user *User) (*InlineCommandProcessorResult, error)
}

// InlineCommandId Идентификатор команды, специально сделан числом, так как приходится эту информацию передавать в телеграм, чтобы
// узнавать какую команду пользователь хочет вызвать. Почему выбран int, а не классный и понятный string. Дело в том,
// что у телеграмма есть ограничение на кол-во данные, которые можно передать в callback_query, ограничение 64 байта.
type InlineCommandId int

// InlineCommandState мета данные команды
type InlineCommandState struct {
	// Идентификатор команды (3 байта)
	Id InlineCommandId
	// Строки (1+ байт)
	Strings []string
}

func (m *InlineCommandState) Encode() string {
	str := strconv.Itoa(int(m.Id)) + "#"
	for i, s := range m.Strings {
		if i != 0 {
			str += ";"
		}

		str += s
	}

	return str
}

// Decode распарсивает данные команды
func (m *InlineCommandState) Decode(str string) error {
	parts := strings.Split(str, "#")
	if len(parts) != 2 {
		return errors.New("invalid format")
	}

	commandIdAsInt, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("can't parse command id: %w", err)
	}

	m.Id = InlineCommandId(commandIdAsInt)
	if len(parts) == 1 || parts[1] == "" {
		return nil
	}

	strs := strings.Split(parts[1], ";")
	if len(strs) == 1 && strs[0] == "" {
		return nil
	}

	m.Strings = append(m.Strings, strs...)

	return nil
}

func NewInlineCommandProcessorResult(
	msgs []tgbotapi.Chattable,
	newPlace *Place,
	commandCall *InlineCommandState,
) *InlineCommandProcessorResult {
	return &InlineCommandProcessorResult{
		msgs:        msgs,
		newPlace:    newPlace,
		commandCall: commandCall,
	}
}

type InlineCommandProcessorResult struct {
	// Список сообщений, которые надо отправить пользователю
	msgs []tgbotapi.Chattable
	// Новое состояние, которое необходимо проставить пользователю
	newPlace *Place
	// Команда, которую надо вызвать
	commandCall *InlineCommandState
}

func (r *InlineCommandProcessorResult) Msgs() []tgbotapi.Chattable {
	return r.msgs
}

func (r *InlineCommandProcessorResult) NewPlace() *Place {
	return r.newPlace
}

func (r *InlineCommandProcessorResult) CommandCall() *InlineCommandState {
	return r.commandCall
}

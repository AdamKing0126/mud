package utils

import (
	"math"
	"mud/commands"
	"mud/notifications"
	worldState "mud/world_state"
	"strings"
)

// type CommandHandlerWithPriority struct {
// 	Handler    CommandHandler
// 	Notifier   notifications.Notifier
// 	WorldState worldState.WorldState
// 	Priority   int
// }

// type CommandHandler interface {
// 	Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string))
// }

type Notifiable interface {
	SetNotifier(notifier *notifications.Notifier)
}

type UsesWorldState interface {
	SetWorldState(worldState *worldState.WorldState)
}

type CommandParser struct {
	commandName string
	arguments   []string
}

func NewCommandParser(commandString string, CommandHandlers map[string]commands.CommandHandlerWithPriority) *CommandParser {
	// Split the command string into the command name and arguments.
	commandParts := strings.Split(commandString, " ")
	commandName := commandParts[0]
	arguments := commandParts[1:]

	// look for partial commands ie "n" for "north"
	bestMatch := ""
	highestPriority := math.MaxInt32
	for fullCommand, handlerWithPriority := range CommandHandlers {
		if strings.HasPrefix(strings.ToLower(fullCommand), commandParts[0]) && handlerWithPriority.Priority < highestPriority {
			bestMatch = fullCommand
			highestPriority = handlerWithPriority.Priority
		}
	}
	if bestMatch != "" {
		commandName = bestMatch
	}

	return &CommandParser{
		commandName: commandName,
		arguments:   arguments,
	}
}

func (p *CommandParser) GetCommandName() string {
	return p.commandName
}

func (p *CommandParser) GetArguments() []string {
	return p.arguments
}

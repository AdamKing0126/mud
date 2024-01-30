package commands

import (
	"math"
	"mud/notifications"
	worldState "mud/world_state"
	"strings"
)

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

func NewCommandParser(commandString string, CommandHandlers map[string]CommandHandlerWithPriority) *CommandParser {
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

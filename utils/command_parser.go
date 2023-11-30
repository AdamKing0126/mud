package utils

import (
	"database/sql"
	"mud/interfaces"
	"strings"
)

type CommandHandler interface {
	Execute(db *sql.DB, player interfaces.PlayerInterface, command string, arguments []string, currentChannel chan interfaces.ActionInterface, updateChannel func(string))
}

type CommandParser struct {
	commandName string
	arguments   []string
}

func NewCommandParser(commandString string, CommandHandlers map[string]CommandHandler) *CommandParser {
	// Split the command string into the command name and arguments.
	commandParts := strings.Split(commandString, " ")
	commandName := commandParts[0]
	arguments := commandParts[1:]
	// look for partial commands ie "n" for "north"
	for fullCommand := range CommandHandlers {
		if strings.HasPrefix(strings.ToLower(fullCommand), commandParts[0]) {
			return &CommandParser{
				commandName: fullCommand,
				arguments:   arguments,
			}
		}
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

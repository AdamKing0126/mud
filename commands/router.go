package commands

import (
	"database/sql"
	"fmt"
	"mud/interfaces"
	"mud/utils"
	"strings"
	"sync"
)

type CommandRouterInterface interface {
	HandleCommand(db *sql.DB, player interfaces.PlayerInterface, command []byte, currentChannel chan interfaces.ActionInterface, updateChannel func(string))
}

type CommandRouter struct {
	Handlers map[string]utils.CommandHandler
	mu       sync.RWMutex
}

func NewCommandRouter() *CommandRouter {
	return &CommandRouter{
		Handlers: make(map[string]utils.CommandHandler),
		mu:       sync.RWMutex{},
	}
}

func (r *CommandRouter) RegisterHandler(command string, handler utils.CommandHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Handlers[command] = handler
}

func RegisterCommands(router *CommandRouter, commands map[string]utils.CommandHandlerWithPriority) {
	for command, handlerWithPriority := range commands {
		router.RegisterHandler(command, handlerWithPriority.Handler)
	}
}

func (r *CommandRouter) HandleCommand(db *sql.DB, player interfaces.PlayerInterface, command []byte, currentChannel chan interfaces.ActionInterface, updateChannel func(string)) {
	playerConn := player.GetConn()
	// Convert the command []byte to a string and trim the extra characters off.
	commandString := strings.ToLower(strings.TrimSpace(string(command)))

	commandBlocks := strings.Split(commandString, ";")
	for _, command := range commandBlocks {
		// Parse the command string.
		commandParser := utils.NewCommandParser(strings.TrimSpace(command), CommandHandlers)

		// Get the command name and arguments.
		commandName := commandParser.GetCommandName()
		arguments := commandParser.GetArguments()

		// Check if the command is registered.
		r.mu.RLock()
		defer r.mu.RUnlock()

		handler, ok := r.Handlers[commandName]
		if !ok {
			fmt.Fprintf(playerConn, "Unknown command: %s\n", command)
			return
		}

		handler.Execute(db, player, command, arguments, currentChannel, updateChannel)
	}
}

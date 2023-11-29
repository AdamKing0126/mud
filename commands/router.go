package commands

import (
	"database/sql"
	"fmt"
	"mud/interfaces"
	"strings"
	"sync"
)

type CommandRouterInterface interface {
	HandleCommand(db *sql.DB, player interfaces.PlayerInterface, command []byte, currentChannel chan interfaces.ActionInterface, updateChannel func(string))
}

type CommandRouter struct {
	Handlers map[string]CommandHandler
	mu       sync.RWMutex
}

func NewCommandRouter() *CommandRouter {
	return &CommandRouter{
		Handlers: make(map[string]CommandHandler),
		mu:       sync.RWMutex{},
	}
}

func (r *CommandRouter) RegisterHandler(command string, handler CommandHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Handlers[command] = handler
}

func RegisterCommands(router *CommandRouter, commands map[string]CommandHandler) {
	for command, handler := range commands {
		router.RegisterHandler(command, handler)
	}
}

func (r *CommandRouter) HandleCommand(db *sql.DB, player interfaces.PlayerInterface, command []byte, currentChannel chan interfaces.ActionInterface, updateChannel func(string)) {
	playerConn := player.GetConn()
	// Convert the command []byte to a string and trim the extra characters off.
	commandString := strings.ToLower(strings.TrimSpace(string(command)))

	commandBlocks := strings.Split(commandString, ";")
	for _, command := range commandBlocks {
		// Parse the command string.
		commandParser := NewCommandParser(strings.TrimSpace(command))

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

		if commandHandler, ok := handler.(CommandHandler); ok {
			commandHandler(db, player, command, arguments, currentChannel, updateChannel)
		} else if function, ok := handler.(func()); ok {
			function()
		}

		// Handle the command.
		// handler(db, player, commandName, arguments, currentChannel, updateChannel)
	}
}

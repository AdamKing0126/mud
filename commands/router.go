package commands

import (
	"database/sql"
	"fmt"
	"mud/interfaces"
	"strings"
	"sync"
)

type CommandRouterInterface interface {
	HandleCommand(db *sql.DB, player interfaces.PlayerInterface, command []byte)
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

func (r *CommandRouter) HandleCommand(db *sql.DB, player interfaces.PlayerInterface, command []byte) {
	playerConn := player.GetConn()
	// Convert the command []byte to a string and trim the extra characters off.
	commandString := strings.ToLower(strings.TrimSpace(string(command)))

	// Parse the command string.
	commandParser := NewCommandParser(commandString)

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

	// Handle the command.
	handler(db, player, commandName, arguments)
}

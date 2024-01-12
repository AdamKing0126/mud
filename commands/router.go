package commands

import (
	"database/sql"
	"fmt"
	"mud/display"
	"mud/interfaces"
	"mud/navigation"
	"mud/notifications"
	"mud/utils"
	"strings"
	"sync"
)

type CommandRouterInterface interface {
	HandleCommand(db *sql.DB, player interfaces.Player, command []byte, currentChannel chan interfaces.Action, updateChannel func(string))
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

func RegisterCommands(router *CommandRouter, notifier *notifications.Notifier, navigator *navigation.Navigator, commands map[string]utils.CommandHandlerWithPriority) {
	for command, handlerWithPriority := range commands {
		if notifiable, ok := handlerWithPriority.Handler.(utils.Notifiable); ok {
			notifiable.SetNotifier(notifier)
		}
		if navigatable, ok := handlerWithPriority.Handler.(utils.Navigatable); ok {
			navigatable.SetNavigator(navigator)
		}
		router.RegisterHandler(command, handlerWithPriority.Handler)
	}
}

func (r *CommandRouter) HandleCommand(db *sql.DB, player interfaces.Player, command []byte, currentChannel chan interfaces.Action, updateChannel func(string)) {
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
			display.PrintWithColor(player, fmt.Sprintf("Unknown command: %s\n", command), "danger")
			return
		}

		handler.Execute(db, player, command, arguments, currentChannel, updateChannel)
	}
}

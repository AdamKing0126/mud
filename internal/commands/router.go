package commands

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/adamking0126/mud/internal/display"
	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/internal/game/players"
	world_state "github.com/adamking0126/mud/internal/game/world_state"
	"github.com/adamking0126/mud/internal/notifications"
	"github.com/adamking0126/mud/pkg/database"

	"github.com/jmoiron/sqlx"
)

type CommandRouterInterface interface {
	HandleCommand(db *sqlx.DB, player players.Player, command []byte, currentChannel chan areas.Action, updateChannel func(string))
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

func RegisterCommands(router *CommandRouter, notifier *notifications.Notifier, worldState *world_state.WorldState, playerService *players.Service, areaService *areas.Service, commands map[string]CommandHandlerWithPriority) {
	for command, handlerWithPriority := range commands {
		if notifiable, ok := handlerWithPriority.Handler.(Notifiable); ok {
			notifiable.SetNotifier(notifier)
		}
		if worldStateable, ok := handlerWithPriority.Handler.(UsesWorldState); ok {
			worldStateable.SetWorldState(worldState)
		}
		if playerServiceable, ok := handlerWithPriority.Handler.(UsesPlayerService); ok {
			playerServiceable.SetPlayerService(playerService)
		}
		if areaServiceable, ok := handlerWithPriority.Handler.(UsesAreaService); ok {
			areaServiceable.SetAreaService(areaService)
		}
		router.RegisterHandler(command, handlerWithPriority.Handler)
	}
}

func (r *CommandRouter) HandleCommand(ctx context.Context, db database.DB, player *players.Player, command []byte, currentChannel chan areas.Action, updateChannel func(string)) {
	// Convert the command []byte to a string and trim the extra characters off.
	commandString := strings.ToLower(strings.TrimSpace(string(command)))

	commandBlocks := strings.Split(commandString, ";")
	for _, command := range commandBlocks {
		// Parse the command string.
		commandParser := NewCommandParser(strings.TrimSpace(command), CommandHandlers)

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

		handler.Execute(ctx, db, player, command, arguments, currentChannel, updateChannel)
	}
}

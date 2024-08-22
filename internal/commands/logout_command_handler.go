package commands

import (
	"context"
	"fmt"

	"github.com/adamking0126/mud/internal/display"
	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/internal/game/players"
	"github.com/adamking0126/mud/internal/game/world_state"
	"github.com/adamking0126/mud/internal/notifications"
	"github.com/adamking0126/mud/pkg/database"
)

type LogoutCommandHandler struct {
	Notifier   *notifications.Notifier
	WorldState *world_state.WorldState
}

func (h *LogoutCommandHandler) Execute(ctx context.Context, db database.DB, player *players.Player, _ string, _ []string, _ chan areas.Action, _ func(string)) {
	display.PrintWithColor(player, "Goodbye!\n", "reset")
	if err := player.Logout(ctx, db); err != nil {
		fmt.Printf("Error logging out player: %v\n", err)
		return
	}

	err := h.WorldState.RemovePlayerFromRoom(player.RoomUUID, player)
	if err != nil {
		fmt.Printf("error removing player %s from room %s - %v", player.UUID, player.RoomUUID, err)
		return
	}

	h.Notifier.NotifyRoom(player.RoomUUID, player.UUID, fmt.Sprintf("\n%s has left the game.\n", player.Name))
}

func (h *LogoutCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

func (h *LogoutCommandHandler) SetWorldState(worldState *world_state.WorldState) {
	h.WorldState = worldState
}

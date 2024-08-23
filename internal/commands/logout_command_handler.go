package commands

import (
	"context"
	"fmt"

	"github.com/adamking0126/mud/internal/display"
	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/internal/game/players"
	"github.com/adamking0126/mud/internal/game/world_state"
	"github.com/adamking0126/mud/internal/notifications"
)

type LogoutCommandHandler struct {
	Notifier          *notifications.Notifier
	WorldStateService *world_state.Service
	PlayerService     *players.Service
}

func (h *LogoutCommandHandler) Execute(ctx context.Context, player *players.Player, _ string, _ []string, _ chan areas.Action, _ func(string)) {
	display.PrintWithColor(player, "Goodbye!\n", "reset")
	if err := h.PlayerService.LogoutPlayer(ctx, player); err != nil {
		fmt.Printf("Error logging out player: %v\n", err)
		return
	}

	err := h.WorldStateService.RemovePlayerFromRoom(ctx, player.RoomUUID, player)
	if err != nil {
		fmt.Printf("error removing player %s from room %s - %v", player.UUID, player.RoomUUID, err)
		return
	}

	h.Notifier.NotifyRoom(player.RoomUUID, player.UUID, fmt.Sprintf("\n%s has left the game.\n", player.Name))
}

func (h *LogoutCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

func (h *LogoutCommandHandler) SetWorldStateService(worldStateService *world_state.Service) {
	h.WorldStateService = worldStateService
}

func (h *LogoutCommandHandler) SetPlayerService(playerService *players.Service) {
	h.PlayerService = playerService
}

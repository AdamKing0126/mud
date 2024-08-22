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

type DropCommandHandler struct {
	Notifier    *notifications.Notifier
	AreaService *areas.Service
	WorldState  *world_state.WorldState
}

func (h *DropCommandHandler) SetWorldState(world_state *world_state.WorldState) {
	h.WorldState = world_state
}

func (h *DropCommandHandler) Execute(ctx context.Context, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	roomUUID := player.RoomUUID
	room := h.WorldState.GetRoom(ctx, roomUUID, false)

	playerItems := player.Inventory
	if len(playerItems) > 0 {
		for _, item := range playerItems {
			if item.GetName() == arguments[0] {
				if err := player.RemoveItem(item); err != nil {
					fmt.Printf("error removing item: %s", err)
				}

				h.AreaService.AddItemToRoom(ctx, room, item)
				display.PrintWithColor(player, fmt.Sprintf("You drop the %s.\n", item.GetName()), "reset")
				h.Notifier.NotifyRoom(player.RoomUUID, player.UUID, fmt.Sprintf("\n%s dropped %s.\n", player.Name, item.Name))
				break
			}
		}
	} else {
		display.PrintWithColor(player, "You don't have that item.\n", "warning")
	}
}

func (h *DropCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

package commands

import (
	"context"
	"fmt"

	"github.com/adamking0126/mud/internal/display"
	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/internal/game/players"
	world_state "github.com/adamking0126/mud/internal/game/world_state"
	"github.com/adamking0126/mud/internal/notifications"

	"github.com/adamking0126/mud/pkg/database"
)

type TakeCommandHandler struct {
	Notifier   *notifications.Notifier
	WorldState *world_state.WorldState
}

func (h *TakeCommandHandler) SetWorldState(world_state *world_state.WorldState) {
	h.WorldState = world_state
}

func (h *TakeCommandHandler) Execute(ctx context.Context, db database.DB, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	roomUUID := player.RoomUUID
	currentRoom := h.WorldState.GetRoom(ctx, roomUUID, false)
	items := currentRoom.Items

	if len(items) > 0 {
		for _, item := range items {
			if item.GetName() == arguments[0] {
				if err := currentRoom.RemoveItem(item); err != nil {
					display.PrintWithColor(player, fmt.Sprintf("error removing item from room: %v", err), "reset")
					break
				}
				player.AddItem(ctx, db, item)

				display.PrintWithColor(player, fmt.Sprintf("You take the %s.\n", item.GetName()), "reset")
				h.Notifier.NotifyRoom(player.RoomUUID, player.UUID, fmt.Sprintf("\n%s takes %s.\n", player.Name, item.Name))
				break
			}
		}
	} else {
		display.PrintWithColor(player, "You don't see that here.\n", "reset")
	}
}

func (h *TakeCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

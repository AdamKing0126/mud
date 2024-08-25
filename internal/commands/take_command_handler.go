package commands

import (
	"context"
	"fmt"

	"github.com/adamking0126/mud/internal/display"
	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/internal/game/players"
	world_state "github.com/adamking0126/mud/internal/game/world_state"
	"github.com/adamking0126/mud/internal/notifications"
)

type TakeCommandHandler struct {
	Notifier          *notifications.Notifier
	WorldStateService *world_state.Service
	PlayerService     *players.Service
}

func (h *TakeCommandHandler) SetWorldStateService(worldStateService *world_state.Service) {
	h.WorldStateService = worldStateService
}

func (h *TakeCommandHandler) SetPlayerService(playerService *players.Service) {
	h.PlayerService = playerService
}

func (h *TakeCommandHandler) Execute(ctx context.Context, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	roomUUID := player.RoomUUID
	currentRoom := h.WorldStateService.GetRoomByUUID(ctx, roomUUID)
	items := currentRoom.Items

	if len(items) > 0 {
		for _, item := range items {
			if item.GetName() == arguments[0] {
				h.WorldStateService.RemoveItemFromRoom(ctx, currentRoom, item)
				if err := currentRoom.RemoveItem(item); err != nil {
					display.PrintWithColor(player, fmt.Sprintf("error removing item from room: %v", err), "reset")
					break
				}
				h.PlayerService.AddItemToPlayerInventory(ctx, player, item)

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

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

type GiveCommandHandler struct {
	Notifier          *notifications.Notifier
	WorldStateService *world_state.Service
	PlayerService     *players.Service
}

func (h *GiveCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

func (h *GiveCommandHandler) SetWorldStateService(worldStateService *world_state.Service) {
	h.WorldStateService = worldStateService
}

func (h *GiveCommandHandler) SetPlayerService(playerService *players.Service) {
	h.PlayerService = playerService
}

func (h *GiveCommandHandler) Execute(ctx context.Context, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	item := player.GetItemFromInventory(arguments[0])
	if item == nil {
		display.PrintWithColor(player, "You don't have that item", "reset")
		return
	}

	currentRoomUUID := player.RoomUUID
	currentRoom := h.WorldStateService.GetRoomByUUID(ctx, currentRoomUUID)

	recipient := currentRoom.GetPlayerByName(arguments[1])
	if recipient == nil {
		display.PrintWithColor(player, "You don't see them here", "reset")
		return
	}

	h.PlayerService.RemoveItemFromPlayer(ctx, player, item)
	h.PlayerService.AddItemToPlayer(ctx, recipient, item)

	display.PrintWithColor(player, fmt.Sprintf("You give %s to %s\n", item.GetName(), recipient.Name), "reset")
	h.Notifier.NotifyPlayer(recipient.UUID, fmt.Sprintf("\n%s gives you %s\n", player.Name, item.GetName()))
}

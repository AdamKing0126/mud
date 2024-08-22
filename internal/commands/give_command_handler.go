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

type GiveCommandHandler struct {
	Notifier   *notifications.Notifier
	WorldState *world_state.WorldState
}

func (h *GiveCommandHandler) SetNotifier(notifier *notifications.Notifier, world_state *world_state.WorldState) {
	h.Notifier = notifier
	h.WorldState = world_state
}

func (h *GiveCommandHandler) Execute(ctx context.Context, db database.DB, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	item := player.GetItemFromInventory(arguments[0])
	if item == nil {
		display.PrintWithColor(player, "You don't have that item", "reset")
		return
	}

	currentRoomUUID := player.RoomUUID
	currentRoom := h.WorldState.GetRoom(ctx, currentRoomUUID, false)

	recipient := currentRoom.GetPlayerByName(arguments[1])
	if recipient == nil {
		display.PrintWithColor(player, "You don't see them here", "reset")
		return
	}

	player.RemoveItem(item)
	recipient.AddItem(ctx, db, item)

	display.PrintWithColor(player, fmt.Sprintf("You give %s to %s\n", item.GetName(), recipient.Name), "reset")
	h.Notifier.NotifyPlayer(recipient.UUID, fmt.Sprintf("\n%s gives you %s\n", player.Name, item.GetName()))
}

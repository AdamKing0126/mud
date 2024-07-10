package commands

import (
	"fmt"
	"mud/areas"
	"mud/display"
	"mud/notifications"
	"mud/players"
	"mud/world_state"

	"github.com/jmoiron/sqlx"
)

type GiveCommandHandler struct {
	Notifier   *notifications.Notifier
	WorldState *world_state.WorldState
}

func (h *GiveCommandHandler) SetNotifier(notifier *notifications.Notifier, world_state *world_state.WorldState) {
	h.Notifier = notifier
	h.WorldState = world_state
}

func (h *GiveCommandHandler) Execute(db *sqlx.DB, player players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	item := player.GetItemFromInventory(arguments[0])
	if item == nil {
		display.PrintWithColor(player, "You don't have that item", "reset")
		return
	}

	currentRoomUUID := player.GetRoomUUID()
	currentRoom := h.WorldState.GetRoom(currentRoomUUID, false)

	recipient := currentRoom.GetPlayerByName(arguments[1])
	if recipient == nil {
		display.PrintWithColor(player, "You don't see them here", "reset")
		return
	}

	player.RemoveItem(*item)
	recipient.AddItem(db, *item)

	display.PrintWithColor(player, fmt.Sprintf("You give %s to %s\n", item.GetName(), recipient.GetName()), "reset")
	h.Notifier.NotifyPlayer(recipient.GetUUID(), fmt.Sprintf("\n%s gives you %s\n", player.GetName(), item.GetName()))
}

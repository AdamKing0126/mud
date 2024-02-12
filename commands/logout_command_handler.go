package commands

import (
	"fmt"
	"mud/display"
	"mud/interfaces"
	"mud/notifications"
	"mud/world_state"

	"github.com/jmoiron/sqlx"
)

type LogoutCommandHandler struct {
	Notifier   *notifications.Notifier
	WorldState *world_state.WorldState
}

func (h *LogoutCommandHandler) Execute(db *sqlx.DB, player interfaces.Player, _ string, _ []string, _ chan interfaces.Action, _ func(string)) {
	display.PrintWithColor(player, "Goodbye!\n", "reset")
	if err := player.Logout(db); err != nil {
		fmt.Printf("Error logging out player: %v\n", err)
		return
	}

	err := h.WorldState.RemovePlayerFromRoom(player.GetRoomUUID(), player)
	if err != nil {
		fmt.Printf("error removing player %s from room %s - %v", player.GetUUID(), player.GetRoomUUID(), err)
		return
	}

	h.Notifier.NotifyRoom(player.GetRoomUUID(), player.GetUUID(), fmt.Sprintf("\n%s has left the game.\n", player.GetName()))
}

func (h *LogoutCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

func (h *LogoutCommandHandler) SetWorldState(worldState *world_state.WorldState) {
	h.WorldState = worldState
}

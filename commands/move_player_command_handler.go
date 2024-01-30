package commands

import (
	"database/sql"
	"fmt"
	"mud/display"
	"mud/interfaces"
	"mud/notifications"
	"mud/world_state"
)

type MovePlayerCommandHandler struct {
	Direction  string
	Notifier   *notifications.Notifier
	WorldState *world_state.WorldState
}

func movePlayerToDirection(worldState *world_state.WorldState, db *sql.DB, player interfaces.Player, room interfaces.Room, direction string, notifier *notifications.Notifier, world_state *world_state.WorldState, currentChannel chan interfaces.Action, updateChannel func(string)) {
	if room == nil || room.GetUUID() == "" {
		display.PrintWithColor(player, "You cannot go that way.", "reset")
	} else {
		display.PrintWithColor(player, "=======================\n\n", "secondary")
		notifier.NotifyRoom(player.GetRoomUUID(), player.GetUUID(), fmt.Sprintf("\n%s goes %s.\n", player.GetName(), direction))

		worldState.RemovePlayerFromRoom(player.GetRoomUUID(), player)
		worldState.AddPlayerToRoom(room.GetUUID(), player)

		notifier.NotifyRoom(room.GetUUID(), player.GetUUID(), fmt.Sprintf("\n%s has arrived.\n", player.GetName()))

		player.SetLocation(db, room.GetUUID())
		var lookArgs []string
		lookHandler := &LookCommandHandler{WorldState: world_state}
		lookHandler.Execute(db, player, "look", lookArgs, currentChannel, updateChannel)
	}
}

func (h *MovePlayerCommandHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	areaUUID := player.GetAreaUUID()

	currentRoom := player.GetRoom()
	exits := currentRoom.GetExits()

	switch h.Direction {
	case "north":
		movePlayerToDirection(h.WorldState, db, player, exits.GetNorth(), h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	case "south":
		movePlayerToDirection(h.WorldState, db, player, exits.GetSouth(), h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	case "west":
		movePlayerToDirection(h.WorldState, db, player, exits.GetWest(), h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	case "east":
		movePlayerToDirection(h.WorldState, db, player, exits.GetEast(), h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	case "up":
		movePlayerToDirection(h.WorldState, db, player, exits.GetUp(), h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	default:
		movePlayerToDirection(h.WorldState, db, player, exits.GetDown(), h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	}

	if areaUUID != player.GetAreaUUID() {
		updateChannel(player.GetAreaUUID())
	}
}

func (h *MovePlayerCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

func (h *MovePlayerCommandHandler) SetWorldState(world_state *world_state.WorldState) {
	h.WorldState = world_state
}

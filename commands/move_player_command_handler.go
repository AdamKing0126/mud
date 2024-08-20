package commands

import (
	"fmt"

	"github.com/adamking0126/mud/areas"
	"github.com/adamking0126/mud/display"
	"github.com/adamking0126/mud/notifications"
	"github.com/adamking0126/mud/players"
	world_state "github.com/adamking0126/mud/world_state"

	"github.com/jmoiron/sqlx"
)

type MovePlayerCommandHandler struct {
	Direction  string
	Notifier   *notifications.Notifier
	WorldState *world_state.WorldState
}

func movePlayerToDirection(worldState *world_state.WorldState, db *sqlx.DB, player *players.Player, room *areas.Room, direction string, notifier *notifications.Notifier, world_state *world_state.WorldState, currentChannel chan areas.Action, updateChannel func(string)) {
	if room == nil || room.UUID == "" {
		display.PrintWithColor(player, "You cannot go that way.", "reset")
	} else {
		display.PrintWithColor(player, "=======================\n\n", "secondary")
		notifier.NotifyRoom(player.RoomUUID, player.UUID, fmt.Sprintf("\n%s goes %s.\n", player.Name, direction))

		worldState.RemovePlayerFromRoom(player.RoomUUID, player)
		worldState.AddPlayerToRoom(room.UUID, player)

		notifier.NotifyRoom(room.UUID, player.UUID, fmt.Sprintf("\n%s has arrived.\n", player.Name))

		player.SetLocation(db, room.UUID)
		var lookArgs []string
		lookHandler := &LookCommandHandler{WorldState: world_state}
		lookHandler.Execute(db, player, "look", lookArgs, currentChannel, updateChannel)
	}
}

func (h *MovePlayerCommandHandler) Execute(db *sqlx.DB, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	areaUUID := player.AreaUUID

	currentRoomUUID := player.RoomUUID
	currentRoom := h.WorldState.GetRoom(currentRoomUUID, true)
	exits := currentRoom.Exits

	switch h.Direction {
	case "north":
		movePlayerToDirection(h.WorldState, db, player, exits.North, h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	case "south":
		movePlayerToDirection(h.WorldState, db, player, exits.South, h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	case "west":
		movePlayerToDirection(h.WorldState, db, player, exits.West, h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	case "east":
		movePlayerToDirection(h.WorldState, db, player, exits.East, h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	case "up":
		movePlayerToDirection(h.WorldState, db, player, exits.Up, h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	default:
		movePlayerToDirection(h.WorldState, db, player, exits.Down, h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	}

	if areaUUID != player.AreaUUID {
		updateChannel(player.AreaUUID)
	}
}

func (h *MovePlayerCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

func (h *MovePlayerCommandHandler) SetWorldState(world_state *world_state.WorldState) {
	h.WorldState = world_state
}

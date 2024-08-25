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

type MovePlayerCommandHandler struct {
	Direction         string
	Notifier          *notifications.Notifier
	WorldStateService *world_state.Service
	PlayerService     *players.Service
}

func movePlayerToDirection(
	ctx context.Context,
	worldStateService *world_state.Service,
	player *players.Player,
	room *areas.Room,
	direction string,
	notifier *notifications.Notifier,
	currentChannel chan areas.Action,
	updateChannel func(string)) {

	if room == nil || room.UUID == "" {
		display.PrintWithColor(player, "You cannot go that way.", "reset")
	} else {
		display.PrintWithColor(player, "=======================\n\n", "secondary")
		notifier.NotifyRoom(player.RoomUUID, player.UUID, fmt.Sprintf("\n%s goes %s.\n", player.Name, direction))

		worldStateService.RemovePlayerFromRoom(ctx, player.RoomUUID, player)
		err := worldStateService.AddPlayerToRoom(ctx, room.UUID, player)
		if err != nil {
			display.PrintWithColor(player, "Error moving player to room: "+err.Error(), "error")
		} else {
			notifier.NotifyRoom(room.UUID, player.UUID, fmt.Sprintf("\n%s has arrived.\n", player.Name))
		}

		var lookArgs []string
		lookHandler := &LookCommandHandler{WorldStateService: worldStateService}
		lookHandler.Execute(ctx, player, "look", lookArgs, currentChannel, updateChannel)
	}
}

func (h *MovePlayerCommandHandler) Execute(ctx context.Context, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	areaUUID := player.AreaUUID

	currentRoomUUID := player.RoomUUID
	currentRoom := h.WorldStateService.GetRoom(ctx, currentRoomUUID, true)
	exits := currentRoom.Exits

	switch h.Direction {
	case "north":
		movePlayerToDirection(ctx, h.WorldStateService, player, exits.North, h.Direction, h.Notifier, currentChannel, updateChannel)
	case "south":
		movePlayerToDirection(ctx, h.WorldStateService, player, exits.South, h.Direction, h.Notifier, currentChannel, updateChannel)
	case "west":
		movePlayerToDirection(ctx, h.WorldStateService, player, exits.West, h.Direction, h.Notifier, currentChannel, updateChannel)
	case "east":
		movePlayerToDirection(ctx, h.WorldStateService, player, exits.East, h.Direction, h.Notifier, currentChannel, updateChannel)
	case "up":
		movePlayerToDirection(ctx, h.WorldStateService, player, exits.Up, h.Direction, h.Notifier, currentChannel, updateChannel)
	default:
		movePlayerToDirection(ctx, h.WorldStateService, player, exits.Down, h.Direction, h.Notifier, currentChannel, updateChannel)
	}

	if areaUUID != player.AreaUUID {
		updateChannel(player.AreaUUID)
	}
}

func (h *MovePlayerCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

func (h *MovePlayerCommandHandler) SetWorldStateService(worldStateService *world_state.Service) {
	h.WorldStateService = worldStateService
}

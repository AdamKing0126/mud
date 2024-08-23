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
	playerService *players.Service,
	player *players.Player,
	room *areas.Room,
	direction string,
	notifier *notifications.Notifier,
	world_state *world_state.WorldState,
	currentChannel chan areas.Action,
	updateChannel func(string)) {

	if room == nil || room.UUID == "" {
		display.PrintWithColor(player, "You cannot go that way.", "reset")
	} else {
		display.PrintWithColor(player, "=======================\n\n", "secondary")
		notifier.NotifyRoom(player.RoomUUID, player.UUID, fmt.Sprintf("\n%s goes %s.\n", player.Name, direction))

		worldStateService.RemovePlayerFromRoom(player.RoomUUID, player)
		worldStateService.AddPlayerToRoom(room.UUID, player)

		notifier.NotifyRoom(room.UUID, player.UUID, fmt.Sprintf("\n%s has arrived.\n", player.Name))

		playerService.SetLocation(ctx, room.UUID)
		var lookArgs []string
		lookHandler := &LookCommandHandler{WorldState: world_state}
		lookHandler.Execute(ctx, worldStateService, playerService, player, "look", lookArgs, currentChannel, updateChannel)
	}
}

func (h *MovePlayerCommandHandler) Execute(ctx context.Context, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	areaUUID := player.AreaUUID

	currentRoomUUID := player.RoomUUID
	currentRoom := h.WorldStateService.GetRoomByUUID(ctx, currentRoomUUID, true)
	exits := currentRoom.Exits

	switch h.Direction {
	case "north":
		movePlayerToDirection(ctx, player, exits.North, h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	case "south":
		movePlayerToDirection(ctx, player, exits.South, h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	case "west":
		movePlayerToDirection(ctx, h.WorldStateService, player, exits.West, h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	case "east":
		movePlayerToDirection(ctx, h.WorldStateService, player, exits.East, h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	case "up":
		movePlayerToDirection(ctx, h.WorldStateService, player, exits.Up, h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	default:
		movePlayerToDirection(ctx, h.WorldStateService, player, exits.Down, h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
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

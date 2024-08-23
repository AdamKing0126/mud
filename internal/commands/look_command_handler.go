package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/adamking0126/mud/internal/display"
	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/internal/game/players"
	"github.com/adamking0126/mud/internal/game/world_state"
)

type LookCommandHandler struct {
	WorldStateService *world_state.Service
}

func (h *LookCommandHandler) SetWorldStateService(world_state_service *world_state.Service) {
	h.WorldStateService = world_state_service
}

func (h *LookCommandHandler) Execute(ctx context.Context, worldStateService *world_state.Service, playerService *players.Service, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	currentRoomUUID := player.RoomUUID
	currentRoom := h.WorldStateService.GetRoom(ctx, currentRoomUUID, false)

	if len(arguments) == 0 {
		display.PrintWithColor(player, fmt.Sprintf("%s\n", currentRoom.Name), "primary")
		display.PrintWithColor(player, fmt.Sprintf("%s\n", currentRoom.Description), "secondary")
		display.PrintWithColor(player, "-----------------------\n\n", "secondary")

		if len(currentRoom.Items) > 0 {
			display.PrintWithColor(player, "You see the following items:\n", "reset")
			for _, item := range currentRoom.Items {
				display.PrintWithColor(player, fmt.Sprintf("%s\n", item.Name), "primary")
			}
			display.PrintWithColor(player, "\n", "reset")
		}

		if len(currentRoom.Mobs) > 0 {
			for _, mob := range currentRoom.Mobs {
				display.PrintWithColor(player, fmt.Sprintf("%s\n", mob.Name), "warning")
			}

		}

		if len(currentRoom.Players) > 1 {
			display.PrintWithColor(player, "You see the following players:\n", "reset")
			for _, playerInRoom := range currentRoom.Players {
				if player.UUID != playerInRoom.UUID {
					display.PrintWithColor(player, fmt.Sprintf("%s\n", playerInRoom.Name), "primary")
				}
			}
			display.PrintWithColor(player, "\n", "reset")
		}

		exitsHandler := &ExitsCommandHandler{ShowOnlyDirections: true, WorldStateService: h.WorldStateService}
		exitsHandler.Execute(ctx, worldStateService, playerService, player, "exits", arguments, currentChannel, updateChannel)
	} else if len(arguments) == 1 {
		exits := currentRoom.Exits
		exitMap := map[string]*areas.Room{
			"North": exits.North,
			"South": exits.South,
			"West":  exits.West,
			"East":  exits.East,
			"Up":    exits.Up,
			"Down":  exits.Down,
		}

		lookDirection := arguments[0]
		directionMatch := false

		for direction, exit := range exitMap {
			if lookDirection == direction {
				directionMatch = true
				if exit != nil {
					exitRoom := h.WorldStateService.GetRoom(ctx, exit.UUID, false)
					display.PrintWithColor(player, fmt.Sprintf("You look %s.  You see %s\n", direction, exitRoom.Name), "reset")
				} else {
					display.PrintWithColor(player, "You don't see anything in that direction\n", "reset")
				}
			}
		}

		if !directionMatch {
			target := arguments[0]
			found := false

			items := append(currentRoom.Items, player.Inventory...)
			for _, item := range items {
				if item.Name == target {
					display.PrintWithColor(player, fmt.Sprintf("%s\n", item.Name), "reset")
					found = true
					break
				}
			}

			for _, playerInRoom := range currentRoom.Players {
				if strings.ToLower(playerInRoom.Name) == target {
					display.PrintWithColor(player, fmt.Sprintf("You see %s.\n", playerInRoom.Name), "reset")
					found = true
					break
				}
			}

			for _, mobInRoom := range currentRoom.Mobs {
				if strings.ToLower(mobInRoom.Name) == target {
					display.PrintWithColor(player, fmt.Sprintf("You see %s.\n", mobInRoom.Name), "danger")
					found = true
					break
				}
			}

			if !found {
				display.PrintWithColor(player, "You don't see that.\n", "reset")
			}
		}
	} else {
		display.PrintWithColor(player, "I don't know how to do that yet.\n", "reset")
	}
}

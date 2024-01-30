package commands

import (
	"database/sql"
	"fmt"
	"mud/display"
	"mud/interfaces"
	"mud/items"
	"mud/world_state"
	"strings"
)

type LookCommandHandler struct {
	WorldState *world_state.WorldState
}

func (h *LookCommandHandler) SetWorldState(world_state *world_state.WorldState) {
	h.WorldState = world_state
}

func (h *LookCommandHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	currentRoom := player.GetRoom()

	if len(arguments) == 0 {
		display.PrintWithColor(player, fmt.Sprintf("%s\n", currentRoom.GetName()), "primary")
		display.PrintWithColor(player, fmt.Sprintf("%s\n", currentRoom.GetDescription()), "secondary")
		display.PrintWithColor(player, "-----------------------\n\n", "secondary")

		if len(currentRoom.GetItems()) > 0 {
			display.PrintWithColor(player, "You see the following items:\n", "reset")
			for _, item := range currentRoom.GetItems() {
				display.PrintWithColor(player, fmt.Sprintf("%s\n", item.GetName()), "primary")
			}
			display.PrintWithColor(player, "\n", "reset")
		}

		if len(currentRoom.GetPlayers()) > 1 {
			display.PrintWithColor(player, "You see the following players:\n", "reset")
			for _, playerInRoom := range currentRoom.GetPlayers() {
				if player.GetUUID() != playerInRoom.GetUUID() {
					display.PrintWithColor(player, fmt.Sprintf("%s\n", playerInRoom.GetName()), "primary")
				}
			}
			display.PrintWithColor(player, "\n", "reset")
		}

		exitsHandler := &ExitsCommandHandler{ShowOnlyDirections: true, WorldState: h.WorldState}
		exitsHandler.Execute(db, player, "exits", arguments, currentChannel, updateChannel)
	} else if len(arguments) == 1 {
		exits := currentRoom.GetExits()
		exitMap := map[string]interfaces.Room{
			"North": exits.GetNorth(),
			"South": exits.GetSouth(),
			"West":  exits.GetWest(),
			"East":  exits.GetEast(),
			"Up":    exits.GetUp(),
			"Down":  exits.GetDown(),
		}

		lookDirection := arguments[0]
		directionMatch := false

		for direction, exit := range exitMap {
			if lookDirection == direction {
				directionMatch = true
				if exit != nil {
					exitRoom := h.WorldState.GetRoom(exit.GetUUID(), false)
					display.PrintWithColor(player, fmt.Sprintf("You look %s.  You see %s\n", direction, exitRoom.GetName()), "reset")
				} else {
					display.PrintWithColor(player, "You don't see anything in that direction\n", "reset")
				}
			}
		}

		if !directionMatch {
			target := arguments[0]
			found := false
			itemsForPlayer, err := items.GetItemsForPlayer(db, player.GetUUID())
			if err != nil {
				display.PrintWithColor(player, fmt.Sprintf("%v", err), "danger")
			}

			items := append(currentRoom.GetItems(), itemsForPlayer...)
			for _, item := range items {
				if item.GetName() == target {
					display.PrintWithColor(player, fmt.Sprintf("%s\n", item.GetName()), "reset")
					found = true
					break
				}
			}

			for _, playerInRoom := range currentRoom.GetPlayers() {
				if strings.ToLower(playerInRoom.GetName()) == target {
					display.PrintWithColor(player, fmt.Sprintf("You see %s.\n", playerInRoom.GetName()), "reset")
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

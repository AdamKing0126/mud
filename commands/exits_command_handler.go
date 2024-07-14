package commands

import (
	"fmt"
	"mud/areas"
	"mud/display"
	"mud/players"
	"mud/world_state"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ExitsCommandHandler struct {
	ShowOnlyDirections bool
	WorldState         *world_state.WorldState
}

func (h *ExitsCommandHandler) SetWorldState(world_state *world_state.WorldState) {
	h.WorldState = world_state
}

func (h *ExitsCommandHandler) Execute(_ *sqlx.DB, player *players.Player, _ string, _ []string, _ chan areas.Action, _ func(string)) {
	roomUUID := player.RoomUUID
	currentRoom := h.WorldState.GetRoom(roomUUID, true)
	exits := currentRoom.Exits
	exitMap := map[string]*areas.Room{
		"North": exits.GetNorth(),
		"South": exits.GetSouth(),
		"West":  exits.GetWest(),
		"East":  exits.GetEast(),
		"Up":    exits.GetUp(),
		"Down":  exits.GetDown(),
	}

	abbreviatedDirections := []string{}
	longDirections := []string{}

	for direction, exit := range exitMap {
		if exit != nil {
			abbreviatedDirections = append(abbreviatedDirections, direction)
			exitRoom := h.WorldState.GetRoom(exit.UUID, false)
			longDirections = append(longDirections, fmt.Sprintf("%s: %s", direction, exitRoom.Name))
		}
	}
	if h.ShowOnlyDirections {
		display.PrintWithColor(player, fmt.Sprintf("\nExits: %s\n", strings.Join(abbreviatedDirections, ", ")), "reset")
	} else {
		for _, direction := range longDirections {
			display.PrintWithColor(player, fmt.Sprintf("%s\n", direction), "reset")
		}
	}
}

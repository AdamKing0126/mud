package commands

import (
	"fmt"

	"github.com/adamking0126/mud/areas"
	"github.com/adamking0126/mud/display"
	"github.com/adamking0126/mud/players"
	"github.com/adamking0126/mud/world_state"

	"github.com/jmoiron/sqlx"
)

type AreaCommandHandler struct {
	WorldState *world_state.WorldState
}

func (h *AreaCommandHandler) SetWorldState(world_state *world_state.WorldState) {
	h.WorldState = world_state
}

func (h *AreaCommandHandler) Execute(db *sqlx.DB, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	areaUUID := player.AreaUUID
	area := h.WorldState.GetArea(areaUUID)

	display.PrintWithColor(player, fmt.Sprintf("%s\n", area.Name), "primary")
	display.PrintWithColor(player, fmt.Sprintf("%s\n", area.Description), "secondary")
	display.PrintWithColor(player, "-----------------------\n\n", "secondary")
}

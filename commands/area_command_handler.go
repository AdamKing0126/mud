package commands

import (
	"fmt"
	"mud/areas"
	"mud/display"
	"mud/players"
	"mud/world_state"

	"github.com/jmoiron/sqlx"
)

type AreaCommandHandler struct {
	WorldState *world_state.WorldState
}

func (h *AreaCommandHandler) SetWorldState(world_state *world_state.WorldState) {
	h.WorldState = world_state
}

func (h *AreaCommandHandler) Execute(db *sqlx.DB, player players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	areaUUID := player.GetAreaUUID()
	area := h.WorldState.GetArea(areaUUID)

	display.PrintWithColor(player, fmt.Sprintf("%s\n", area.GetName()), "primary")
	display.PrintWithColor(player, fmt.Sprintf("%s\n", area.GetDescription()), "secondary")
	display.PrintWithColor(player, "-----------------------\n\n", "secondary")
}

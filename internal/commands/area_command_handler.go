package commands

import (
	"context"
	"fmt"

	"github.com/adamking0126/mud/internal/display"
	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/internal/game/players"
	"github.com/adamking0126/mud/internal/game/world_state"
	"github.com/adamking0126/mud/pkg/database"
)

type AreaCommandHandler struct {
	WorldState *world_state.WorldState
}

func (h *AreaCommandHandler) SetWorldState(world_state *world_state.WorldState) {
	h.WorldState = world_state
}

func (h *AreaCommandHandler) Execute(ctx context.Context, db database.DB, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	areaUUID := player.AreaUUID
	area := h.WorldState.GetArea(ctx, areaUUID)

	display.PrintWithColor(player, fmt.Sprintf("%s\n", area.Name), "primary")
	display.PrintWithColor(player, fmt.Sprintf("%s\n", area.Description), "secondary")
	display.PrintWithColor(player, "-----------------------\n\n", "secondary")
}

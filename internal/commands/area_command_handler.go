package commands

import (
	"context"
	"fmt"

	"github.com/adamking0126/mud/internal/display"
	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/internal/game/players"
	"github.com/adamking0126/mud/internal/game/world_state"
)

type AreaCommandHandler struct {
	WorldStateService *world_state.Service
}

func (h *AreaCommandHandler) SetWorldStateService(worldStateService *world_state.Service) {
	h.WorldStateService = worldStateService
}

func (h *AreaCommandHandler) Execute(ctx context.Context, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	areaUUID := player.AreaUUID
	area := h.WorldStateService.GetAreaByUUID(ctx, areaUUID)

	display.PrintWithColor(player, fmt.Sprintf("%s\n", area.Name), "primary")
	display.PrintWithColor(player, fmt.Sprintf("%s\n", area.Description), "secondary")
	display.PrintWithColor(player, "-----------------------\n\n", "secondary")
}

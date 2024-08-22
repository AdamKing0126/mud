package commands

import (
	"context"
	"fmt"

	"github.com/adamking0126/mud/internal/display"
	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/internal/game/players"
	"github.com/adamking0126/mud/pkg/database"
)

type InventoryCommandHandler struct{}

func (h *InventoryCommandHandler) Execute(ctx context.Context, db database.DB, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	display.PrintWithColor(player, "You are carrying:\n", "secondary")
	playerInventory := player.Inventory

	if len(playerInventory) == 0 {
		display.PrintWithColor(player, "Nothing\n", "reset")
	} else {
		for _, item := range playerInventory {
			display.PrintWithColor(player, fmt.Sprintf("%s\n", item.Name), "reset")
		}
	}
}

package commands

import (
	"fmt"

	"github.com/adamking0126/mud/areas"
	"github.com/adamking0126/mud/display"
	"github.com/adamking0126/mud/players"

	"github.com/jmoiron/sqlx"
)

type InventoryCommandHandler struct{}

func (h *InventoryCommandHandler) Execute(db *sqlx.DB, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
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

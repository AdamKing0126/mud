package commands

import (
	"database/sql"
	"fmt"
	"mud/display"
	"mud/interfaces"
)

type InventoryCommandHandler struct{}

func (h *InventoryCommandHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	display.PrintWithColor(player, "You are carrying:\n", "secondary")
	playerInventory := player.GetInventory()

	if len(playerInventory) == 0 {
		display.PrintWithColor(player, "Nothing\n", "reset")
	} else {
		for _, item := range playerInventory {
			display.PrintWithColor(player, fmt.Sprintf("%s\n", item.GetName()), "reset")
		}
	}
}

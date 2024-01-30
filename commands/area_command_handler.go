package commands

import (
	"database/sql"
	"fmt"
	"mud/display"
	"mud/interfaces"
)

type AreaCommandHandler struct{}

func (h *AreaCommandHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	area := player.GetArea()
	display.PrintWithColor(player, fmt.Sprintf("%s\n", area.GetName()), "primary")
	display.PrintWithColor(player, fmt.Sprintf("%s\n", area.GetDescription()), "secondary")
	display.PrintWithColor(player, "-----------------------\n\n", "secondary")
}

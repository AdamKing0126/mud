package commands

import (
	"database/sql"
	"mud/display"
	"mud/interfaces"
)

type WhoAmICommandHandler struct{}

func (*WhoAmICommandHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	display.PrintWithColor(player, player.GetName(), "reset")

}

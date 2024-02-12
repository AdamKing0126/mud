package commands

import (
	"mud/display"
	"mud/interfaces"

	"github.com/jmoiron/sqlx"
)

type WhoAmICommandHandler struct{}

func (*WhoAmICommandHandler) Execute(db *sqlx.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	display.PrintWithColor(player, player.GetName(), "reset")

}

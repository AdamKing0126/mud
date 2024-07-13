package commands

import (
	"mud/areas"
	"mud/display"
	"mud/players"

	"github.com/jmoiron/sqlx"
)

type WhoAmICommandHandler struct{}

func (*WhoAmICommandHandler) Execute(db *sqlx.DB, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	display.PrintWithColor(player, player.Name, "reset")

}

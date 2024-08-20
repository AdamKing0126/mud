package commands

import (
	"github.com/adamking0126/mud/internal/display"
	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/internal/game/players"

	"github.com/jmoiron/sqlx"
)

type WhoAmICommandHandler struct{}

func (*WhoAmICommandHandler) Execute(db *sqlx.DB, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	display.PrintWithColor(player, player.Name, "reset")

}

package commands

import (
	"mud/areas"
	"mud/display"
	"mud/notifications"
	"mud/players"

	"github.com/jmoiron/sqlx"
)

type RemoveCommandHandler struct {
	Notifier *notifications.Notifier
}

func (h *RemoveCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

func (h *RemoveCommandHandler) Execute(db *sqlx.DB, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	if len(arguments) == 0 {
		display.PrintWithColor(player, "Remove what?\n", "primary")
		return
	}

	player.Remove(db, arguments[0])

}

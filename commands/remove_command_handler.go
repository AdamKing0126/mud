package commands

import (
	"database/sql"
	"mud/display"
	"mud/interfaces"
	"mud/notifications"
)

type RemoveCommandHandler struct {
	Notifier *notifications.Notifier
}

func (h *RemoveCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

func (h *RemoveCommandHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	if len(arguments) == 0 {
		display.PrintWithColor(player, "Remove what?\n", "primary")
		return
	}

	player.Remove(db, arguments[0])

}

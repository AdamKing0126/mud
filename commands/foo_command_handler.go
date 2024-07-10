package commands

import (
	"mud/areas"
	"mud/notifications"
	"mud/players"

	"github.com/jmoiron/sqlx"
)

type FooCommandHandler struct {
	Notifier *notifications.Notifier
}

func (h *FooCommandHandler) Execute(db *sqlx.DB, player players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	currentChannel <- areas.Action{Player: player, Command: command, Arguments: arguments}
}

func (h *FooCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

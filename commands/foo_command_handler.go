package commands

import (
	"mud/areas"
	"mud/interfaces"
	"mud/notifications"

	"github.com/jmoiron/sqlx"
)

type FooCommandHandler struct {
	Notifier *notifications.Notifier
}

func (h *FooCommandHandler) Execute(db *sqlx.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	currentChannel <- &areas.Action{Player: player, Command: command, Arguments: arguments}
}

func (h *FooCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

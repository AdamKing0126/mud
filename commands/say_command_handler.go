package commands

import (
	"fmt"
	"mud/areas"
	"mud/display"
	"mud/notifications"
	"mud/players"
	"strings"

	"github.com/jmoiron/sqlx"
)

type SayHandler struct {
	Notifier *notifications.Notifier
}

func (h *SayHandler) Execute(db *sqlx.DB, player players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	msg := strings.Join(arguments, " ")
	display.PrintWithColor(player, fmt.Sprintf("You say \"%s\"\n", msg), "reset")
	h.Notifier.NotifyRoom(player.GetRoomUUID(), player.GetUUID(), fmt.Sprintf("\n%s says \"%s\"\n", player.GetName(), msg))
}

func (h *SayHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

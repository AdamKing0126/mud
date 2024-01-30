package commands

import (
	"database/sql"
	"fmt"
	"mud/display"
	"mud/interfaces"
	"mud/notifications"
	"strings"
)

type SayHandler struct {
	Notifier *notifications.Notifier
}

func (h *SayHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	msg := strings.Join(arguments, " ")
	display.PrintWithColor(player, fmt.Sprintf("You say \"%s\"\n", msg), "reset")
	h.Notifier.NotifyRoom(player.GetRoomUUID(), player.GetUUID(), fmt.Sprintf("\n%s says \"%s\"\n", player.GetName(), msg))
}

func (h *SayHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

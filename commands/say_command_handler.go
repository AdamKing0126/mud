package commands

import (
	"fmt"
	"strings"

	"github.com/adamking0126/mud/areas"
	"github.com/adamking0126/mud/display"
	"github.com/adamking0126/mud/notifications"
	"github.com/adamking0126/mud/players"

	"github.com/jmoiron/sqlx"
)

type SayHandler struct {
	Notifier *notifications.Notifier
}

func (h *SayHandler) Execute(db *sqlx.DB, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	msg := strings.Join(arguments, " ")
	display.PrintWithColor(player, fmt.Sprintf("You say \"%s\"\n", msg), "reset")
	h.Notifier.NotifyRoom(player.RoomUUID, player.UUID, fmt.Sprintf("\n%s says \"%s\"\n", player.Name, msg))
}

func (h *SayHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

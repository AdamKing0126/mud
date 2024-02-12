package commands

import (
	"fmt"
	"mud/display"
	"mud/interfaces"
	"mud/notifications"
	"mud/players"
	"strings"

	"github.com/jmoiron/sqlx"
)

type TellHandler struct {
	Notifier *notifications.Notifier
}

func (h *TellHandler) Execute(db *sqlx.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	msg := strings.Join(arguments[1:], " ")
	retrievedPlayer, err := players.GetPlayerByName(db, arguments[0])
	if err != nil {
		display.PrintWithColor(player, fmt.Sprintf("Error retrieving player UUID: %v\n", err), "danger")
		return
	}

	if player.GetUUID() == retrievedPlayer.GetUUID() {
		display.PrintWithColor(player, "Talking to yourself again?\n", "reset")
		return
	}

	if retrievedPlayer.GetLoggedIn() {
		display.PrintWithColor(player, fmt.Sprintf("You tell %s \"%s\"\n", arguments[0], msg), "reset")
		h.Notifier.NotifyPlayer(retrievedPlayer.GetUUID(), fmt.Sprintf("\n%s tells you \"%s\"\n", player.GetName(), msg))
	} else {
		display.PrintWithColor(player, fmt.Sprintf("%s isn't here\n", arguments[0]), "reset")
	}
}

func (h *TellHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

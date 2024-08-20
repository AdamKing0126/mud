package commands

import (
	"fmt"
	"strings"

	"github.com/adamking0126/mud/internal/game/areas"

	"github.com/adamking0126/mud/internal/display"
	"github.com/adamking0126/mud/internal/game/players"
	"github.com/adamking0126/mud/internal/notifications"

	"github.com/jmoiron/sqlx"
)

type TellHandler struct {
	Notifier *notifications.Notifier
}

func (h *TellHandler) Execute(db *sqlx.DB, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	msg := strings.Join(arguments[1:], " ")
	retrievedPlayer, err := players.GetPlayerByName(db, arguments[0])
	if err != nil {
		display.PrintWithColor(player, fmt.Sprintf("Error retrieving player UUID: %v\n", err), "danger")
		return
	}

	if player.UUID == retrievedPlayer.UUID {
		display.PrintWithColor(player, "Talking to yourself again?\n", "reset")
		return
	}

	if retrievedPlayer.LoggedIn {
		display.PrintWithColor(player, fmt.Sprintf("You tell %s \"%s\"\n", arguments[0], msg), "reset")
		h.Notifier.NotifyPlayer(retrievedPlayer.UUID, fmt.Sprintf("\n%s tells you \"%s\"\n", player.Name, msg))
	} else {
		display.PrintWithColor(player, fmt.Sprintf("%s isn't here\n", arguments[0]), "reset")
	}
}

func (h *TellHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

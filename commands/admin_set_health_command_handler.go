package commands

import (
	"fmt"
	"mud/areas"
	"mud/display"
	"mud/notifications"
	"mud/players"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type AdminSetHealthCommandHandler struct {
	Notifier *notifications.Notifier
}

func (h *AdminSetHealthCommandHandler) Execute(db *sqlx.DB, player players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	target := arguments[0]
	value := arguments[1]

	retrievedPlayer, err := players.GetPlayerByName(db, target)
	if err != nil {
		display.PrintWithColor(player, fmt.Sprintf("Error retrieving player UUID: %v\n", err), "danger")
		return
	}

	query := "UPDATE players SET health = ? WHERE UUID = ?"

	_, err = db.Exec(query, value, retrievedPlayer.GetUUID())

	if err != nil {
		display.PrintWithColor(player, fmt.Sprintf("Error updating health: %v\n", err), "danger")
		return
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		display.PrintWithColor(player, fmt.Sprintf("Error converting value to int: %v\n", err), "danger")
		return
	}
	h.Notifier.Players[retrievedPlayer.GetUUID()].SetHealth(int32(intValue))
	display.PrintWithColor(player, fmt.Sprintf("You set %s's health to %d\n", target, intValue), "reset")
	h.Notifier.NotifyPlayer(retrievedPlayer.GetUUID(), fmt.Sprintf("\n%s magically sets your health to %d\n", player.GetName(), intValue))

}

func (h *AdminSetHealthCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

package commands

import (
	"fmt"
	"mud/display"
	"mud/interfaces"
	"mud/notifications"

	"github.com/jmoiron/sqlx"
)

type GiveCommandHandler struct {
	Notifier *notifications.Notifier
}

func (h *GiveCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

func (h *GiveCommandHandler) Execute(db *sqlx.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	item := player.GetItemFromInventory(arguments[0])
	if item == nil {
		display.PrintWithColor(player, "You don't have that item", "reset")
		return
	}

	currentRoom := player.GetRoom()
	recipient := currentRoom.GetPlayerByName(arguments[1])
	if recipient == nil {
		display.PrintWithColor(player, "You don't see them here", "reset")
		return
	}

	player.RemoveItem(item)
	recipient.AddItem(db, item)

	display.PrintWithColor(player, fmt.Sprintf("You give %s to %s\n", item.GetName(), recipient.GetName()), "reset")
	h.Notifier.NotifyPlayer(recipient.GetUUID(), fmt.Sprintf("\n%s gives you %s\n", player.GetName(), item.GetName()))
}

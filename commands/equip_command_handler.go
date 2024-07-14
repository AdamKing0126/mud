package commands

import (
	"fmt"
	"mud/areas"
	"mud/display"
	"mud/notifications"
	"mud/players"

	"github.com/jmoiron/sqlx"
)

type EquipHandler struct {
	Notifier *notifications.Notifier
}

func (h *EquipHandler) Execute(db *sqlx.DB, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	if len(arguments) == 0 {
		player.DisplayEquipment()
		return
	}

	playerItems := player.Inventory

	if len(playerItems) > 0 {
		for _, item := range playerItems {
			if item.GetName() == arguments[0] {
				if player.Equip(db, item) {
					display.PrintWithColor(player, fmt.Sprintf("You wield %s.\n", item.GetName()), "reset")
					h.Notifier.NotifyRoom(player.RoomUUID, player.UUID, fmt.Sprintf("\n%s wields %s.\n", player.Name, item.Name))
					player.RemoveItem(item)
				}
				break
			}
		}
	} else {
		display.PrintWithColor(player, "You don't have that item.\n", "warning")
	}
}

func (h *EquipHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

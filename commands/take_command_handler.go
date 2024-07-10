package commands

import (
	"fmt"
	"mud/areas"
	"mud/display"
	"mud/notifications"
	"mud/players"

	"github.com/jmoiron/sqlx"
)

type TakeCommandHandler struct {
	Notifier *notifications.Notifier
}

func (h *TakeCommandHandler) Execute(db *sqlx.DB, player players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	currentRoom := player.GetRoom()
	items := currentRoom.GetItems()

	if len(items) > 0 {
		for _, item := range items {
			if item.GetName() == arguments[0] {
				if err := currentRoom.RemoveItem(item); err != nil {
					display.PrintWithColor(player, fmt.Sprintf("error removing item from room: %v", err), "reset")
					break
				}
				player.AddItem(db, item)

				display.PrintWithColor(player, fmt.Sprintf("You take the %s.\n", item.GetName()), "reset")
				h.Notifier.NotifyRoom(player.GetRoomUUID(), player.GetUUID(), fmt.Sprintf("\n%s takes %s.\n", player.GetName(), item.GetName()))
				break
			}
		}
	} else {
		display.PrintWithColor(player, "You don't see that here.\n", "reset")
	}
}

func (h *TakeCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

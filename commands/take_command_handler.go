package commands

import (
	"database/sql"
	"fmt"
	"mud/display"
	"mud/interfaces"
	"mud/notifications"
)

type TakeCommandHandler struct {
	Notifier *notifications.Notifier
}

func (h *TakeCommandHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
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

package commands

import (
	"fmt"
	"mud/areas"
	"mud/display"
	"mud/notifications"
	"mud/players"

	"github.com/jmoiron/sqlx"
)

type DropCommandHandler struct {
	Notifier *notifications.Notifier
}

func (h *DropCommandHandler) Execute(db *sqlx.DB, player players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	roomUUID := player.GetRoomUUID()
	room := player.GetRoom()

	playerItems := player.GetInventory()
	if len(playerItems) > 0 {
		for _, item := range playerItems {
			if item.GetName() == arguments[0] {
				if err := player.RemoveItem(item); err != nil {
					fmt.Printf("error removing item: %s", err)
				}
				room.AddItem(db, item)
				query := "UPDATE item_locations SET room_uuid = ?, player_uuid = NULL WHERE item_uuid = ?"
				_, err := db.Exec(query, roomUUID, item.GetUUID())
				if err != nil {
					display.PrintWithColor(player, fmt.Sprintf("Failed to update item location: %v\n", err), "danger")
				}
				display.PrintWithColor(player, fmt.Sprintf("You drop the %s.\n", item.GetName()), "reset")
				h.Notifier.NotifyRoom(player.GetRoomUUID(), player.GetUUID(), fmt.Sprintf("\n%s dropped %s.\n", player.GetName(), item.GetName()))
				break
			}
		}
	} else {
		display.PrintWithColor(player, "You don't have that item.\n", "warning")
	}
}

func (h *DropCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

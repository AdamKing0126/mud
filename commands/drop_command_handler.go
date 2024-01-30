package commands

import (
	"database/sql"
	"fmt"
	"mud/display"
	"mud/interfaces"
	"mud/notifications"
)

type DropCommandHandler struct {
	Notifier *notifications.Notifier
}

func (h *DropCommandHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
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

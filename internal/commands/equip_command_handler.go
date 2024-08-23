package commands

import (
	"context"
	"fmt"

	"github.com/adamking0126/mud/internal/display"
	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/internal/game/players"
	"github.com/adamking0126/mud/internal/notifications"
)

type EquipHandler struct {
	Notifier      *notifications.Notifier
	PlayerService *players.Service
}

func (h *EquipHandler) Execute(ctx context.Context, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	if len(arguments) == 0 {
		player.DisplayEquipment()
		return
	}

	playerItems := player.Inventory

	if len(playerItems) > 0 {
		for _, item := range playerItems {
			if item.GetName() == arguments[0] {
				if h.PlayerService.EquipItem(ctx, player, item) {
					display.PrintWithColor(player, fmt.Sprintf("You wield %s.\n", item.GetName()), "reset")
					h.Notifier.NotifyRoom(player.RoomUUID, player.UUID, fmt.Sprintf("\n%s wields %s.\n", player.Name, item.Name))
					h.PlayerService.RemoveItemFromPlayer(ctx, player, item)
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

func (h *EquipHandler) SetPlayerService(playerService *players.Service) {
	h.PlayerService = playerService
}

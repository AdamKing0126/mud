package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/adamking0126/mud/internal/display"
	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/internal/game/players"
	"github.com/adamking0126/mud/internal/notifications"
)

type AdminSetHealthCommandHandler struct {
	Notifier      *notifications.Notifier
	PlayerService *players.Service
}

func (h *AdminSetHealthCommandHandler) Execute(ctx context.Context, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	name := arguments[0]
	value := arguments[1]

	targetPlayer, err := h.PlayerService.GetPlayerByName(ctx, name)
	if err != nil {
		display.PrintWithColor(player, fmt.Sprintf("Player %s not found\n", name), "danger")
		return
	}

	convertedValue, err := strconv.Atoi(value)
	if err != nil {
		display.PrintWithColor(player, fmt.Sprintf("Error converting value to int: %v\n", err), "danger")
		return
	}

	h.PlayerService.SetPlayerHealth(ctx, targetPlayer, convertedValue)

	playerInNotifier := h.Notifier.Players[targetPlayer.UUID]
	// TODO what if the player isn't found
	playerInNotifier.HP = int32(convertedValue)
	display.PrintWithColor(player, fmt.Sprintf("You set %s's health to %d\n", name, convertedValue), "reset")
	h.Notifier.NotifyPlayer(targetPlayer.UUID, fmt.Sprintf("\n%s magically sets your health to %d\n", player.Name, convertedValue))

}

func (h *AdminSetHealthCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

func (h *AdminSetHealthCommandHandler) SetPlayerService(playerService *players.Service) {
	h.PlayerService = playerService
}

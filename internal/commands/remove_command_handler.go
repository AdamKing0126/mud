package commands

import (
	"context"

	"github.com/adamking0126/mud/internal/display"
	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/internal/game/players"
	"github.com/adamking0126/mud/internal/notifications"
)

type RemoveCommandHandler struct {
	Notifier      *notifications.Notifier
	PlayerService *players.Service
}

func (h *RemoveCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

func (h *RemoveCommandHandler) SetPlayerService(playerService *players.Service) {
	h.PlayerService = playerService
}

func (h *RemoveCommandHandler) Execute(ctx context.Context, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	if len(arguments) == 0 {
		display.PrintWithColor(player, "Remove what?\n", "primary")
		return
	}

	h.PlayerService.UnequipItem(ctx, player, arguments[0])
}

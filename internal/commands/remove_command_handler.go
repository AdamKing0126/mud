package commands

import (
	"context"

	"github.com/adamking0126/mud/internal/display"
	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/internal/game/players"
	"github.com/adamking0126/mud/internal/notifications"
	"github.com/adamking0126/mud/pkg/database"
)

type RemoveCommandHandler struct {
	Notifier *notifications.Notifier
}

func (h *RemoveCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

func (h *RemoveCommandHandler) Execute(ctx context.Context, db database.DB, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	if len(arguments) == 0 {
		display.PrintWithColor(player, "Remove what?\n", "primary")
		return
	}

	player.Remove(ctx, db, arguments[0])
}

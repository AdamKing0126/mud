package commands

import (
	"context"

	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/internal/game/players"
	"github.com/adamking0126/mud/internal/notifications"
	"github.com/adamking0126/mud/pkg/database"
)

type FooCommandHandler struct {
	Notifier *notifications.Notifier
}

func (h *FooCommandHandler) Execute(ctx context.Context, db database.DB, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	currentChannel <- areas.Action{Player: *player, Command: command, Arguments: arguments}
}

func (h *FooCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

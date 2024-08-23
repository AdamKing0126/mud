package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/adamking0126/mud/internal/display"
	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/internal/game/players"
	worldState "github.com/adamking0126/mud/internal/game/world_state"
	"github.com/adamking0126/mud/internal/notifications"
)

type SayHandler struct {
	Notifier *notifications.Notifier
}

func (h *SayHandler) Execute(ctx context.Context, worldStateService *worldState.Service, playerService *players.Service, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	msg := strings.Join(arguments, " ")
	display.PrintWithColor(player, fmt.Sprintf("You say \"%s\"\n", msg), "reset")
	h.Notifier.NotifyRoom(player.RoomUUID, player.UUID, fmt.Sprintf("\n%s says \"%s\"\n", player.Name, msg))
}

func (h *SayHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

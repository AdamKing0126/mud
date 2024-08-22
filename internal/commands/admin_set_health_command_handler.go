package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/adamking0126/mud/internal/display"
	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/internal/game/players"
	"github.com/adamking0126/mud/internal/notifications"
	"github.com/adamking0126/mud/pkg/database"
)

type AdminSetHealthCommandHandler struct {
	Notifier *notifications.Notifier
}

func (h *AdminSetHealthCommandHandler) Execute(ctx context.Context, db database.DB, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	target := arguments[0]
	value := arguments[1]

	retrievedPlayer, err := players.GetPlayerByName(ctx, db, target)
	if err != nil {
		display.PrintWithColor(player, fmt.Sprintf("Error retrieving player UUID: %v\n", err), "danger")
		return
	}

	query := "UPDATE players SET health = ? WHERE UUID = ?"

	err = db.Exec(ctx, query, value, retrievedPlayer.UUID)

	if err != nil {
		display.PrintWithColor(player, fmt.Sprintf("Error updating health: %v\n", err), "danger")
		return
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		display.PrintWithColor(player, fmt.Sprintf("Error converting value to int: %v\n", err), "danger")
		return
	}
	playerInNotifier := h.Notifier.Players[retrievedPlayer.UUID]
	// TODO what if the player isn't found
	playerInNotifier.HP = int32(intValue)
	display.PrintWithColor(player, fmt.Sprintf("You set %s's health to %d\n", target, intValue), "reset")
	h.Notifier.NotifyPlayer(retrievedPlayer.UUID, fmt.Sprintf("\n%s magically sets your health to %d\n", player.Name, intValue))

}

func (h *AdminSetHealthCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

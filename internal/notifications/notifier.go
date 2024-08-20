package notifications

import (
	"fmt"

	"github.com/adamking0126/mud/internal/display"
	"github.com/adamking0126/mud/internal/game/players"
)

type Notifier struct {
	Players map[string]*players.Player
}

func NewNotifier(connections map[string]*players.Player) *Notifier {
	return &Notifier{Players: connections}
}

func (n *Notifier) NotifyRoom(roomID string, playerUUID string, message string) {
	fmt.Println("Notifying room", roomID, "with message", message)
	var playersInRoom []*players.Player
	for _, player := range n.Players {
		if player.RoomUUID == roomID && player.UUID != playerUUID {
			playersInRoom = append(playersInRoom, player)
		}
	}
	for _, player := range playersInRoom {
		display.PrintWithColor(player, message, "primary")
		display.PrintWithColor(player, fmt.Sprintf("\nHP: %d Mvt: %d> ", player.HP, player.Movement), "primary")
	}
}

func (n *Notifier) NotifyAll(message string) {
	for _, player := range n.Players {
		display.PrintWithColor(player, message, "primary")
		display.PrintWithColor(player, fmt.Sprintf("\nHP: %d Mvt: %d> ", player.HP, player.Movement), "primary")
	}
}

func (n *Notifier) NotifyPlayer(playerUUID string, message string) {
	player := n.Players[playerUUID]
	display.PrintWithColor(player, message, "primary")
	display.PrintWithColor(player, fmt.Sprintf("\nHP: %d Mvt: %d> ", player.HP, player.Movement), "primary")
}

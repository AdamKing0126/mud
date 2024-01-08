package notifications

import (
	"fmt"
	"mud/display"
	"mud/interfaces"
)

type Notifier struct {
	Players map[string]interfaces.Player
}

func NewNotifier(connections map[string]interfaces.Player) *Notifier {
	return &Notifier{Players: connections}
}

func (n *Notifier) NotifyRoom(roomID string, playerUUID string, message string) {
	fmt.Println("Notifying room", roomID, "with message", message)
	var playersInRoom []interfaces.Player
	for _, player := range n.Players {
		if player.GetRoom() == roomID && player.GetUUID() != playerUUID {
			playersInRoom = append(playersInRoom, player)
		}
	}
	for _, player := range playersInRoom {
		display.PrintWithColor(player, message, "primary")
		display.PrintWithColor(player, fmt.Sprintf("\nHP: %d Mana: %d Mvt: %d> ", player.GetHealth(), player.GetMana(), player.GetMovement()), "primary")
	}
}

func (n *Notifier) NotifyAll(message string) {
	for _, player := range n.Players {
		display.PrintWithColor(player, message, "primary")
		display.PrintWithColor(player, fmt.Sprintf("\nHP: %d Mana: %d Mvt: %d> ", player.GetHealth(), player.GetMana(), player.GetMovement()), "primary")
	}
}

func (n *Notifier) NotifyPlayer(playerUUID string, message string) {
	player := n.Players[playerUUID]
	display.PrintWithColor(player, message, "primary")
	display.PrintWithColor(player, fmt.Sprintf("\nHP: %d Mana: %d Mvt: %d> ", player.GetHealth(), player.GetMana(), player.GetMovement()), "primary")
}

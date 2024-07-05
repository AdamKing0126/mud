package commands

import (
	"fmt"
	"mud/combat"
	"mud/display"
	"mud/interfaces"
	"mud/players"

	"github.com/jmoiron/sqlx"
)

type PlayerStatusCommandHandler struct{}

func (h *PlayerStatusCommandHandler) Execute(db *sqlx.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	playerAbilities := &players.PlayerAbilities{}

	query := "SELECT * FROM player_abilities WHERE player_uuid = ?"
	err := db.QueryRow(query, player.GetUUID()).Scan(&playerAbilities.UUID, &playerAbilities.PlayerUUID, &playerAbilities.Strength, &playerAbilities.Intelligence, &playerAbilities.Wisdom, &playerAbilities.Constitution, &playerAbilities.Charisma, &playerAbilities.Dexterity)
	if err != nil {
		display.PrintWithColor(player, fmt.Sprintf("%v", err), "danger")
	}

	player.SetAbilities(playerAbilities)

	display.PrintWithColor(player, fmt.Sprintf("%s\n", player.GetCharacterClass()), "danger")
	display.PrintWithColor(player, fmt.Sprintf("%s\n", player.GetRace()), "danger")
	display.PrintWithColor(player, fmt.Sprintf("Strength: %d\n", playerAbilities.GetStrength()), "danger")
	display.PrintWithColor(player, fmt.Sprintf("Dexterity: %d\n", playerAbilities.GetDexterity()), "danger")
	display.PrintWithColor(player, fmt.Sprintf("Constitution: %d\n", playerAbilities.GetConstitution()), "danger")
	display.PrintWithColor(player, fmt.Sprintf("Intelligence: %d\n", playerAbilities.GetIntelligence()), "danger")
	display.PrintWithColor(player, fmt.Sprintf("Wisdom: %d\n", playerAbilities.GetWisdom()), "danger")
	display.PrintWithColor(player, fmt.Sprintf("Charisma: %d\n", playerAbilities.GetCharisma()), "danger")

	// for debugging purposes only - remove later
	display.PrintWithColor(player, "\n\n***********DEBUG***************\n", "danger")
	display.PrintWithColor(player, fmt.Sprintf("Attack Roll Hits: %t\n", combat.AttackRoll(player, player)), "danger")
	display.PrintWithColor(player, "*******************************\n", "danger")
}

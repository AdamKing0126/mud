package commands

import (
	"context"
	"fmt"

	"github.com/adamking0126/mud/internal/display"
	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/internal/game/players"
	"github.com/adamking0126/mud/internal/game/world_state"
)

type PlayerStatusCommandHandler struct {
	WorldStateService *world_state.Service
	PlayerService     *players.Service
}

func (h *PlayerStatusCommandHandler) SetWorldStateService(worldStateService *world_state.Service) {
	h.WorldStateService = worldStateService
}

func (h *PlayerStatusCommandHandler) SetPlayerService(playerService *players.Service) {
	h.PlayerService = playerService
}

func (h *PlayerStatusCommandHandler) Execute(ctx context.Context, player *players.Player, command string, arguments []string, currentChannel chan areas.Action, updateChannel func(string)) {
	playerAbilities := h.PlayerService.GetPlayerAbilities(ctx, player)

	display.PrintWithColor(player, fmt.Sprintf("%s\n", player.GetCharacterClass()), "danger")
	display.PrintWithColor(player, fmt.Sprintf("%s\n", player.GetRace()), "danger")
	display.PrintWithColor(player, fmt.Sprintf("Strength: %d\n", playerAbilities.GetStrength()), "danger")
	display.PrintWithColor(player, fmt.Sprintf("Dexterity: %d\n", playerAbilities.GetDexterity()), "danger")
	display.PrintWithColor(player, fmt.Sprintf("Constitution: %d\n", playerAbilities.GetConstitution()), "danger")
	display.PrintWithColor(player, fmt.Sprintf("Intelligence: %d\n", playerAbilities.GetIntelligence()), "danger")
	display.PrintWithColor(player, fmt.Sprintf("Wisdom: %d\n", playerAbilities.GetWisdom()), "danger")
	display.PrintWithColor(player, fmt.Sprintf("Charisma: %d\n", playerAbilities.GetCharisma()), "danger")

	// TODO for debugging purposes only - remove later
	// display.PrintWithColor(player, "\n\n***********DEBUG***************\n", "danger")
	// display.PrintWithColor(player, fmt.Sprintf("Attack Roll Hits: %t\n", combat.AttackRoll(player, player)), "danger")
	// display.PrintWithColor(player, "*******************************\n", "danger")
}

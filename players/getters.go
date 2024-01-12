package players

import (
	"mud/interfaces"
	"net"
)

func (player *Player) GetAbilities() interfaces.Abilities {
	return &player.PlayerAbilities
}

func (player *Player) GetArmorClass() int {
	// 10 + armor_bonus + shield_bonus + dexterity_modifier + other_modifiers
	base := 10
	armorBonus := 0
	shieldBonus := 0
	dexModifier := player.PlayerAbilities.GetDexterityModifier()
	otherModifiers := 0
	return base + armorBonus + shieldBonus + dexModifier + otherModifiers
}

func (player *Player) GetArea() string {
	return player.Area
}

func (player *Player) GetColorProfile() interfaces.ColorProfile {
	return &player.ColorProfile
}

func (player *Player) GetCommands() []string {
	return player.Commands
}

func (player *Player) GetConn() net.Conn {
	return player.Conn
}

func (player *Player) GetEquipment() interfaces.PlayerEquipment {
	return &player.Equipment
}

func (player *Player) GetHashedPassword() string {
	return player.Password
}

func (player *Player) GetHealth() int {
	return player.Health
}

func (player *Player) GetHealthMax() int {
	return player.HealthMax
}

func (player *Player) GetInventory() []interfaces.Item {
	return player.Inventory
}

func (player *Player) GetLoggedIn() bool {
	return player.LoggedIn
}

func (player *Player) GetMana() int {
	return player.Mana
}

func (player *Player) GetManaMax() int {
	return player.ManaMax
}

func (player *Player) GetMovement() int {
	return player.Movement
}

func (player *Player) GetMovementMax() int {
	return player.MovementMax
}

func (player *Player) GetName() string {
	return player.Name
}

func (player *Player) GetRoomUUID() string {
	return player.Room
}

func (player *Player) GetSizeModifier() int {
	// Need to update this.  Probably need to move this out, so it can be used by players and monsters

	sizeTable := map[string]int{
		"colossal":   -8,
		"gargantuan": -4,
		"huge":       -2,
		"large":      -1,
		"medium":     0,
		"small":      1,
		"tiny":       2,
	}
	// TODO don't hard code this
	return sizeTable["medium"]
}

func (player *Player) GetUUID() string {
	return player.UUID
}

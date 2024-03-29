package players

import (
	"mud/interfaces"
	"net"
)

func (player *Player) GetAbilities() interfaces.Abilities {
	return &player.PlayerAbilities
}

func (player *Player) GetArmorClass() int32 {
	// 10 + armor_bonus + shield_bonus + dexterity_modifier + other_modifiers
	base := int32(10)
	armorBonus := int32(0)
	shieldBonus := int32(0)
	dexModifier := player.PlayerAbilities.GetDexterityModifier()
	otherModifiers := int32(0)
	return base + armorBonus + shieldBonus + dexModifier + otherModifiers
}

func (player *Player) GetAreaUUID() string {
	return player.AreaUUID
}

func (player *Player) GetArea() interfaces.Area {
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

func (player *Player) GetHealth() int32 {
	return player.Health
}

func (player *Player) GetHealthMax() int32 {
	return player.HealthMax
}

func (player *Player) GetInventory() []interfaces.Item {
	return player.Inventory
}

func (player *Player) GetItemFromInventory(itemName string) interfaces.Item {
	inventory := player.GetInventory()
	for idx := range inventory {
		if inventory[idx].GetName() == itemName {
			return inventory[idx]
		}
	}
	return nil
}

func (player *Player) GetLoggedIn() bool {
	return player.LoggedIn
}

func (player *Player) GetMana() int32 {
	return player.Mana
}

func (player *Player) GetManaMax() int32 {
	return player.ManaMax
}

func (player *Player) GetMovement() int32 {
	return player.Movement
}

func (player *Player) GetMovementMax() int32 {
	return player.MovementMax
}

func (player *Player) GetName() string {
	return player.Name
}

func (player *Player) GetRoomUUID() string {
	return player.RoomUUID
}

func (player *Player) GetRoom() interfaces.Room {
	return player.Room
}

func (player *Player) GetSizeModifier() int32 {
	// Need to update this.  Probably need to move this out, so it can be used by players and monsters

	sizeTable := map[string]int32{
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

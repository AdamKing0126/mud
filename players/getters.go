package players

import (
	"github.com/adamking0126/mud/items"
)

func (player *Player) GetArmorClass() int32 {
	// 10 + armor_bonus + shield_bonus + dexterity_modifier + other_modifiers
	base := int32(10)
	armorBonus := int32(0)
	shieldBonus := int32(0)
	dexModifier := player.PlayerAbilities.GetDexterityModifier()
	otherModifiers := int32(0)
	return base + armorBonus + shieldBonus + dexModifier + otherModifiers
}

func (player *Player) GetCharacterClass() string {
	return player.CharacterClass.Name + " - " + player.CharacterClass.ArchetypeName
}

func (player *Player) GetRace() string {
	if player.Race.SubRaceName == "" {
		return player.Race.Name
	}
	return player.Race.Name + " - " + player.Race.SubRaceName
}

func (player Player) GetItemFromInventory(itemName string) *items.Item {
	inventory := player.Inventory
	for idx := range inventory {
		if inventory[idx].GetName() == itemName {
			return inventory[idx]
		}
	}
	return nil
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

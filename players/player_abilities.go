package players

type PlayerAbilities struct {
	UUID         string
	PlayerUUID   string
	Strength     int32
	Intelligence int32
	Charisma     int32
	Wisdom       int32
	Dexterity    int32
	Constitution int32
}

func (playerAbilities *PlayerAbilities) GetUUID() string {
	return playerAbilities.UUID
}

func (playerAbilities *PlayerAbilities) GetPlayerUUID() string {
	return playerAbilities.PlayerUUID
}

func (playerAbilities *PlayerAbilities) GetStrength() int32 {
	return playerAbilities.Strength
}

func (playerAbilities *PlayerAbilities) GetIntelligence() int32 {
	return playerAbilities.Intelligence
}

func (playerAbilities *PlayerAbilities) GetCharisma() int32 {
	return playerAbilities.Charisma
}

func (playerAbilities *PlayerAbilities) GetWisdom() int32 {
	return playerAbilities.Wisdom
}

func (playerAbilities *PlayerAbilities) GetDexterity() int32 {
	return playerAbilities.Dexterity
}

func (playerAbilities *PlayerAbilities) GetConstitution() int32 {
	return playerAbilities.Constitution
}

func (playerAbilities *PlayerAbilities) GetAttackModifier(weaponType string) int32 {
	if weaponType == "ranged" {
		return playerAbilities.GetDexterityModifier()
	} else if weaponType == "melee" {
		return playerAbilities.GetStrengthModifier()
	}
	return 0
}

func (playerAbilities *PlayerAbilities) GetStrengthModifier() int32 {
	strengthBonusTable := map[int32]int32{
		0:  -6,
		1:  -5,
		2:  -4,
		3:  -4,
		4:  -3,
		5:  -3,
		6:  -2,
		7:  -2,
		8:  -1,
		9:  -1,
		10: 0,
		11: 0,
		12: 1,
		13: 1,
		14: 2,
		15: 2,
		16: 3,
		17: 3,
		18: 4,
		19: 4,
		20: 5,
		21: 5,
		22: 6,
		23: 6,
		24: 7,
		25: 7,
		26: 8,
		27: 8,
		28: 9,
		29: 9,
		30: 10,
		31: 10,
		32: 11,
		33: 11,
	}
	return strengthBonusTable[playerAbilities.GetStrength()]
}

func (playerAbilities *PlayerAbilities) GetDexterityModifier() int32 {
	dexterityBonusTable := map[int32]int32{
		0:  -6,
		1:  -5,
		2:  -4,
		3:  -4,
		4:  -3,
		5:  -3,
		6:  -2,
		7:  -2,
		8:  -1,
		9:  -1,
		10: 0,
		11: 0,
		12: 1,
		13: 1,
		14: 2,
		15: 2,
		16: 3,
		17: 3,
		18: 4,
		19: 4,
		20: 5,
		21: 5,
		22: 6,
		23: 6,
		24: 7,
		25: 7,
		26: 8,
		27: 8,
		28: 9,
		29: 9,
		30: 10,
		31: 10,
		32: 11,
		33: 11,
	}
	return dexterityBonusTable[playerAbilities.GetDexterity()]
}

package mobs

type Mob struct {
	UUID        string
	Name        string
	Description string
	Room        string
	Area        string
	Tags        []string
	Health      int
	HealthMax   int
	Mana        int
	ManaMax     int
	Abilities   MobAbilities
}

func (mob *Mob) GetArmorClass() int {
	return 10
}

type MobAbilities struct{}

func (mobAbilities *MobAbilities) GetAttackModifier(weaponType string) int {
	if weaponType == "ranged" {
		return mobAbilities.GetDexterityModifier()
	} else if weaponType == "melee" {
		return mobAbilities.GetStrengthModifier()
	}
	return 0
}

func (mobAbilities *MobAbilities) GetDexterityModifier() int {
	return 10
}

func (MobAbilities *MobAbilities) GetStrengthModifier() int {
	return 10
}

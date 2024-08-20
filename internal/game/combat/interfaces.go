package combat

type Abilities interface {
	GetAttackModifier(string) int32
	GetStrengthModifier() int32
	GetDexterityModifier() int32
}

type Combatant interface {
	GetAbilities() Abilities
	GetArmorClass() int32
	RollInitiative() int32
}

type Opponent interface {
	GetArmorClass() int32
	GetName() string
}

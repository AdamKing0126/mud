package interfaces

type Combatant interface {
	GetAbilities() Abilities
	GetArmorClass() int32
	RollInitiative() int32
}

package interfaces

type CombatantInterface interface {
	GetAbilities() AbilitiesInterface
	GetArmorClass() int
}

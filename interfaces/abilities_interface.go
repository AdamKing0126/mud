package interfaces

type AbilitiesInterface interface {
	GetAttackModifier(string) int
	GetStrengthModifier() int
	GetDexterityModifier() int
}

package interfaces

type Abilities interface {
	GetAttackModifier(string) int
	GetStrengthModifier() int
	GetDexterityModifier() int
}

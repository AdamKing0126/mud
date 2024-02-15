package interfaces

type Abilities interface {
	GetAttackModifier(string) int32
	GetStrengthModifier() int32
	GetDexterityModifier() int32
}

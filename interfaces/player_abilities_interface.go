package interfaces

type PlayerAbilitiesInterface interface {
	GetStrength() int
	GetIntelligence() int
	GetDexterity() int
	GetCharisma() int
	GetConstitution() int
	GetWisdom() int
	GetPlayerUUID() string
	GetUUID() string
	GetAttackModifier(string) int
	GetStrengthModifier() int
	GetDexterityModifier() int
}

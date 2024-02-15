package interfaces

type PlayerAbilities interface {
	GetStrength() int32
	GetIntelligence() int32
	GetDexterity() int32
	GetCharisma() int32
	GetConstitution() int32
	GetWisdom() int32
	GetPlayerUUID() string
	GetUUID() string
	GetAttackModifier(string) int32
	GetStrengthModifier() int32
	GetDexterityModifier() int32
}

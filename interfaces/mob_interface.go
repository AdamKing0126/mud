package interfaces

type Mob interface {
	GetAreaUUID() string
	GetRoomUUID() string
	GetAlignment() string
	GetArmorClass() int32
	GetArmorDescription() string
	GetChallengeRating() float64
	GetCharisma() int32
	GetCharismaSave() int32
	GetConditionImmunities() string
	GetConstitution() int32
	GetConstitutionSave() int32
	GetDamageImmunities() string
	GetDamageResistances() string
	GetDamageVulnerabilities() string
	GetDescription() string
	GetDexterity() int32
	GetDexteritySave() int32
	GetGroup() string
	GetHP() int32
	GetMaxHP() int32
	GetHitDice() string
	GetIntelligence() int32
	GetIntelligenceSave() int32
	GetLegendaryDescription() string
	GetName() string
	GetPerception() int32
	GetSenses() string
	GetSize() string
	GetSlug() string
	GetStrength() int32
	GetStrengthSave() int32
	GetSubtype() string
	GetType() string
	GetWisdom() int32
	GetWisdomSave() int32
	RollHitDice() int32
	RollInitiative() int32
	GetActions() []MobAction
	ExecuteAction(Opponent)
}

type MobAction interface {
	GetName() string
	GetDescription() string
	GetAttackBonus() int32
	GetDamageDice() string
	GetDamageBonus() int32
	SetDescription(string)
}

package mobs

import "mud/interfaces"

func (mob *Mob) GetAreaUUID() string {
	return mob.AreaUUID
}

func (mob *Mob) GetRoomUUID() string {
	return mob.RoomUUID
}

func (mob *Mob) GetAlignment() string {
	return mob.Alignment
}

func (mob *Mob) GetArmorClass() int32 {
	return mob.ArmorClass
}

func (mob *Mob) GetArmorDescription() string {
	return mob.ArmorDescription
}

func (mob *Mob) GetChallengeRating() float64 {
	return mob.ChallengeRating
}

func (mob *Mob) GetCharisma() int32 {
	return mob.Charisma
}

func (mob *Mob) GetCharismaSave() int32 {
	return mob.CharismaSave
}

func (mob *Mob) GetConditionImmunities() string {
	return mob.ConditionImmunities
}

func (mob *Mob) GetConstitution() int32 {
	return mob.Constitution
}

func (mob *Mob) GetConstitutionSave() int32 {
	return mob.ConstitutionSave
}

func (mob *Mob) GetDamageImmunities() string {
	return mob.DamageImmunities
}

func (mob *Mob) GetDamageResistances() string {
	return mob.DamageResistances
}

func (mob *Mob) GetDamageVulnerabilities() string {
	return mob.DamageVulnerabilities
}

func (mob *Mob) GetDescription() string {
	return mob.Description
}

func (mob *Mob) GetDexterity() int32 {
	return mob.Dexterity
}

func (mob *Mob) GetDexteritySave() int32 {
	return mob.DexteritySave
}

func (mob *Mob) GetGroup() string {
	return mob.Group
}

func (mob *Mob) GetHP() int32 {
	return mob.HP
}

func (mob *Mob) GetMaxHP() int32 {
	return mob.MaxHP
}

func (mob *Mob) GetHitDice() string {
	return mob.HitDice
}

func (mob *Mob) GetIntelligence() int32 {
	return mob.Intelligence
}

func (mob *Mob) GetIntelligenceSave() int32 {
	return mob.IntelligenceSave
}

func (mob *Mob) GetLegendaryDescription() string {
	return mob.LegendaryDescription
}

func (mob *Mob) GetName() string {
	return mob.Name
}

func (mob *Mob) GetPerception() int32 {
	return mob.Perception
}

func (mob *Mob) GetSenses() string {
	return mob.Senses
}

func (mob *Mob) GetSize() string {
	return mob.Size
}

func (mob *Mob) GetSlug() string {
	return mob.Slug
}

func (mob *Mob) GetStrength() int32 {
	return mob.Strength
}

func (mob *Mob) GetStrengthSave() int32 {
	return mob.StrengthSave
}

func (mob *Mob) GetSubtype() string {
	return mob.Subtype
}

func (mob *Mob) GetType() string {
	return mob.Type
}

func (mob *Mob) GetWisdom() int32 {
	return mob.Wisdom
}

func (mob *Mob) GetWisdomSave() int32 {
	return mob.WisdomSave
}

func (mob *Mob) GetActions() []interfaces.MobAction {
	return mob.Actions
}

func (mobAction *Action) GetAttackBonus() int32 {
	return mobAction.AttackBonus
}

func (mobAction *Action) GetDamageBonus() int32 {
	return mobAction.DamageBonus
}

func (mobAction *Action) GetName() string {
	return mobAction.Name
}

func (mobAction *Action) GetDescription() string {
	return mobAction.Description
}

func (mobAction *Action) GetDamageDice() string {
	return mobAction.DamageDice
}

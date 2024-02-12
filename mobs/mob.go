package mobs

import (
	"mud/dice"

	_ "github.com/mattn/go-sqlite3"
)

type Mob struct {
	ID                    int64   `db:"id" mapstructure:"db"`
	AreaUUID              string  `db:"area_uuid" mapstructure:"area_uuid"`
	RoomUUID              string  `db:"room_uuid" mapstructure:"room_uuid"`
	Alignment             string  `db:"alignment" mapstructure:"alignment"`
	ArmorClass            int32   `db:"armor_class" mapstructure:"armor_class"`
	ArmorDescription      string  `db:"armor_description" mapstructure:"armor_description"`
	ChallengeRating       float64 `db:"challenge_rating" mapstructure:"challenge_rating"`
	Charisma              int32   `db:"charisma" mapstructure:"charisma"`
	CharismaSave          int32   `db:"charisma_save" mapstructure:"charisma_save"`
	ConditionImmunities   string  `db:"condition_immunities" mapstructure:"condition_immunities"`
	Constitution          int32   `db:"constitution" mapstructure:"constitution"`
	ConstitutionSave      int32   `db:"constitution_save" mapstructure:"constitution_save"`
	DamageImmunities      string  `db:"damage_immunities" mapstructure:"damage_immunities"`
	DamageResistances     string  `db:"damage_resistances" mapstructure:"damage_resistances"`
	DamageVulnerabilities string  `db:"damage_vulnerabilities" mapstructure:"damage_vulnerabilities"`
	Description           string  `db:"description" mapstructure:"description"`
	Dexterity             int32   `db:"dexterity" mapstructure:"dexterity"`
	DexteritySave         int32   `db:"dexterity_save" mapstructure:"dexterity_save"`
	Group                 string  `db:"group_name" mapstructure:"group_name"`
	HP                    int32   `db:"hp" mapstructure:"hp"`
	MaxHP                 int32   `db:"hp" mapstructure:"hp"`
	HitDice               string  `db:"hit_dice" mapstructure:"hit_dice"`
	Intelligence          int32   `db:"intelligence" mapstructure:"intelligence"`
	IntelligenceSave      int32   `db:"intelligence_save" mapstructure:"intelligence_save"`
	LegendaryDescription  string  `db:"legendary_description" mapstructure:"legendary_description"`
	Name                  string  `db:"name" mapstructure:"name"`
	Perception            int32   `db:"perception" mapstructure:"perception"`
	Senses                string  `db:"senses" mapstructure:"senses"`
	Size                  string  `db:"size" mapstructure:"size"`
	Slug                  string  `db:"slug" mapstructure:"slug"`
	Strength              int32   `db:"strength" mapstructure:"strength"`
	StrengthSave          int32   `db:"strength_save" mapstructure:"strength_save"`
	Subtype               string  `db:"subtype" mapstructure:"subtype"`
	Type                  string  `db:"type" mapstructure:"type"`
	Wisdom                int32   `db:"wisdom" mapstructure:"wisdom"`
	WisdomSave            int32   `db:"wisdom_save" mapstructure:"wisdom_save"`
}

func (mob *Mob) RollHitDice() int {
	return dice.DiceRoll(mob.HitDice)

}

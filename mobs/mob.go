package mobs

import (
	_ "github.com/mattn/go-sqlite3"
)

type Mob struct {
	ID                    int32   `db:"id"`
	AreaUUID              string  `db:"area_uuid"`
	RoomUUID              string  `db:"room_uuid"`
	Alignment             string  `db:"alignment"`
	ArmorClass            int32   `db:"armor_class"`
	ArmorDescription      string  `db:"armor_description"`
	ChallengeRating       float32 `db:"challenge_rating"`
	Charisma              int32   `db:"charisma"`
	CharismaSave          int32   `db:"charisma_save"`
	ConditionImmunities   string  `db:"condition_immunities"`
	Constitution          int32   `db:"constitution"`
	ConstitutionSave      int32   `db:"constitution_save"`
	DamageImmunities      string  `db:"damage_immunities"`
	DamageResistances     string  `db:"damage_resistances"`
	DamageVulnerabilities string  `db:"damage_vulnerabilities"`
	Description           string  `db:"description"`
	Dexterity             int32   `db:"dexterity"`
	DexteritySave         int32   `db:"dexterity_save"`
	Group                 string  `db:"group_name"`
	HP                    int32   `db:"hp"`
	MaxHP                 int32   `db:"hp"`
	HitDice               string  `db:"hit_dice"`
	Intelligence          int32   `db:"intelligence"`
	IntelligenceSave      int32   `db:"intelligence_save"`
	LegendaryDescription  string  `db:"legendary_description"`
	Name                  string  `db:"name"`
	Perception            int32   `db:"perception"`
	Senses                string  `db:"senses"`
	Size                  string  `db:"size"`
	Slug                  string  `db:"slug"`
	Strength              int32   `db:"strength"`
	StrengthSave          int32   `db:"strength_save"`
	Subtype               string  `db:"subtype"`
	Type                  string  `db:"type"`
	Wisdom                int32   `db:"wisdom"`
	WisdomSave            int32   `db:"wisdom_save"`
}

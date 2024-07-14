package mobs

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/mitchellh/mapstructure"
)

// this object represents how mobs are stored in the database.  the fields should
// map to a Mob.
type MobDB struct {
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
	Actions               string  `db:"actions" mapstructure:"actions"`
}

func GetMobsInRoom(db *sqlx.DB, roomUUID string) ([]*Mob, error) {
	var mobs []*Mob
	rows, err := db.Queryx("SELECT * FROM mobs WHERE room_uuid = ?", roomUUID)
	if err != nil {
		log.Fatalf("failed to fetch mobs from tabel: %v", err)
	}

	for rows.Next() {
		result := make(map[string]interface{})
		err = rows.MapScan(result)
		if err != nil {
			log.Fatalf("failed to scan row: %v", err)
		}

		var mobDb MobDB
		err = mapstructure.Decode(result, &mobDb)
		if err != nil {
			log.Fatalf("failed to decode row: %v", err)
		}

		fmt.Printf("Actions: %s\n", mobDb.Actions)
		// TODO I have no idea what I meant to do here...
		fmt.Printf("Type of Action: foo\n")

		var actions []Action
		err = json.Unmarshal([]byte(mobDb.Actions), &actions)
		if err != nil {
			log.Fatalf("error unmarshaling actions: %v", err)
		}

		var mobActions []*Action
		for idx := range actions {
			mobActions = append(mobActions, &actions[idx])
		}

		mob := Mob{
			ID:                    mobDb.ID,
			AreaUUID:              mobDb.AreaUUID,
			RoomUUID:              mobDb.RoomUUID,
			Alignment:             mobDb.Alignment,
			ArmorClass:            mobDb.ArmorClass,
			ArmorDescription:      mobDb.ArmorDescription,
			ChallengeRating:       mobDb.ChallengeRating,
			Charisma:              mobDb.Charisma,
			CharismaSave:          mobDb.CharismaSave,
			ConditionImmunities:   mobDb.ConditionImmunities,
			Constitution:          mobDb.Constitution,
			ConstitutionSave:      mobDb.ConstitutionSave,
			DamageImmunities:      mobDb.DamageImmunities,
			DamageResistances:     mobDb.DamageResistances,
			DamageVulnerabilities: mobDb.DamageVulnerabilities,
			Description:           mobDb.Description,
			Dexterity:             mobDb.Dexterity,
			DexteritySave:         mobDb.DexteritySave,
			Group:                 mobDb.Group,
			HP:                    mobDb.HP,
			MaxHP:                 mobDb.MaxHP,
			HitDice:               mobDb.HitDice,
			Intelligence:          mobDb.Intelligence,
			IntelligenceSave:      mobDb.IntelligenceSave,
			LegendaryDescription:  mobDb.LegendaryDescription,
			Name:                  mobDb.Name,
			Perception:            mobDb.Perception,
			Senses:                mobDb.Senses,
			Size:                  mobDb.Size,
			Slug:                  mobDb.Slug,
			Strength:              mobDb.Strength,
			StrengthSave:          mobDb.StrengthSave,
			Subtype:               mobDb.Subtype,
			Type:                  mobDb.Type,
			Wisdom:                mobDb.Wisdom,
			WisdomSave:            mobDb.WisdomSave,
			// TODO figure out what to do here.
			// Actions:               mobDb.Actions,
		}

		mobs = append(mobs, &mob)
	}

	return mobs, nil

}

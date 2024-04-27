package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type ClassArchetypeImport struct {
	Description        string `json:"desc"`
	DocumentLicenseUrl string `json:"document__license_url"`
	DocumentSlug       string `json:"document__slug"`
	DocumentTitle      string `json:"document__title"`
	DocumentUrl        string `json:"document__url"`
	Name               string `json:"name"`
	Slug               string `json:"slug"`
}

type ClassImport struct {
	Name                      string                 `json:"name" db:"name"`
	Slug                      string                 `json:"slug" db:"slug"`
	Description               string                 `json:"desc" db:"description"`
	HitDice                   string                 `json:"hit_dice" db:"hit_dice"`
	HpAtFirstLevel            string                 `json:"hp_at_1st_level" db:"hp_at_first_level"`
	HpAtHigherLevels          string                 `json:"hp_at_higher_levels" db:"hp_at_higher_levels"`
	ProficienciesArmor        string                 `json:"prof_armor" db:"proficiencies_armor"`
	ProficienciesWeapons      string                 `json:"prof_weapons" db:"proficiencies_weapons"`
	ProficienciesTools        string                 `json:"prof_tools" db:"proficiencies_tools"`
	ProficienciesSavingThrows string                 `json:"prof_saving_throws" db:"proficiencies_saving_throws"`
	ProficienciesSkills       string                 `json:"prof_skills" db:"proficiencies_skills"`
	Equipment                 string                 `json:"equipment" db:"equipment"`
	Table                     string                 `json:"table" db:"class_table"`
	SpellcastingAbility       string                 `json:"spellcasting_ability" db:"spellcasting_ability"`
	SubtypesName              string                 `json:"subtypes_name" db:"subtypes_name"`
	ArchetypesData            string                 `db:"archetypes"`
	Archetypes                []ClassArchetypeImport `json:"archetypes"`
	DocumentSlug              string                 `json:"document__slug" db:"document_slug"`
	DocumentTitle             string                 `json:"document__title" db:"document_title"`
	DocumentLicenseUrl        string                 `json:"document__license_url" db:"document_license_url"`
	DocumentUrl               string                 `json:"document__url" db:"document_url"`
}

func getHpAtFirstLevel(value interface{}) int {
	hpAtFirstLevelStr, ok := value.(string)
	if !ok {
		log.Fatalf("error asserting hp_at_first_level to string")
	}
	hpAtFirstLevel, err := strconv.Atoi(strings.Fields(hpAtFirstLevelStr)[0])
	if err != nil {
		fmt.Println("Error converting string to int: ", err)
	}
	return hpAtFirstLevel
}

func SeedClasses(dbPath string, classesImportDbPath string) error {
	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open Sqlite database: %v", err)
	} else {
		err := db.Ping()
		if err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}
		fmt.Println("Mud database opened successfully")
	}
	defer db.Close()

	classesDB, err := sqlx.Connect("sqlite3", classesImportDbPath)
	if err != nil {
		log.Fatalf("Failed to open Class Imports database: %v", err)
	} else {
		err := db.Ping()
		if err != nil {
			log.Fatalf("Failed to ping database Class Imports: %v", err)
		}
		fmt.Println("Class Imports Database opened successfully")
	}
	defer classesDB.Close()
	query := `SELECT description, hit_dice, hp_at_first_level, hp_at_higher_levels, name, slug, proficiencies_saving_throws, archetypes from class_imports;`
	rows, err := classesDB.Queryx(query)
	if err != nil {
		log.Fatalf("Failed to query row: %v", err)
	}

	for rows.Next() {
		var ci ClassImport
		err = rows.StructScan(&ci)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		var archetypes []ClassArchetypeImport
		err = json.Unmarshal([]byte(ci.ArchetypesData), &archetypes)
		if err != nil {
			log.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		ci.Archetypes = archetypes

		savingThrowValues := strings.Split(ci.ProficienciesSavingThrows, ", ")
		savingThrowCharisma := false
		savingThrowConstitution := false
		savingThrowDexterity := false
		savingThrowIntelligence := false
		savingThrowStrength := false
		savingThrowWisdom := false
		for _, savingThrow := range savingThrowValues {
			savingThrow = strings.ToLower(savingThrow)
			switch savingThrow {
			case "charisma":
				savingThrowCharisma = true
			case "constitution":
				savingThrowConstitution = true
			case "dexterity":
				savingThrowDexterity = true
			case "intelligence":
				savingThrowIntelligence = true
			case "strength":
				savingThrowStrength = true
			case "wisdom":
				savingThrowWisdom = true
			}
		}

		hpModifier := "none"
		for _, ability := range []string{"charisma", "constitution", "dexterity", "intelligence", "strength", "wisdom"} {
			ci.HpAtHigherLevels = strings.ToLower(ci.HpAtHigherLevels)

			if strings.Contains(ci.HpAtHigherLevels, ability) {
				hpModifier = ability
			}
		}

		hpAtFirstLevel := getHpAtFirstLevel(ci.HpAtFirstLevel)

		queryString := `INSERT INTO character_classes 
		(hit_dice, hp_at_first_level, hp_modifier, name, saving_throw_charisma, saving_throw_constitution, saving_throw_dexterity, saving_throw_intelligence, saving_throw_strength, saving_throw_wisdom, slug, archetype_slug, archetype_name, archetype_description)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

		for _, archetype := range ci.Archetypes {
			archetypeDescription := strings.Split(archetype.Description, "#####")[0]
			archetypeDescription = strings.Split(archetypeDescription, "**")[0]
			archetypeDescription = strings.TrimSpace(archetypeDescription)
			_, err = db.Exec(queryString, ci.HitDice, hpAtFirstLevel, hpModifier, ci.Name, savingThrowCharisma, savingThrowConstitution, savingThrowDexterity, savingThrowIntelligence, savingThrowStrength, savingThrowWisdom, ci.Slug, archetype.Slug, archetype.Name, archetypeDescription)
			if err != nil {
				log.Fatalf("Failed to insert SQLite row: %v", err)
			}
		}
	}

	return nil
}

func main() {
	dbPath := "./sql_database/mud.db"
	classesImportDbPath := "./sql_database/class_imports.db"
	SeedClasses(dbPath, classesImportDbPath)
}

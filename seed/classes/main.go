package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

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
		result := make(map[string]interface{})
		err = rows.MapScan(result)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		proficienciesSavingThrows, ok := result["proficiencies_saving_throws"].(string)
		if !ok {
			log.Fatalf("proficienciesSavingThrows is not a string")
		}
		savingThrowValues := strings.Split(proficienciesSavingThrows, ", ")
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
			hpAtHigherLevels := result["hp_at_higher_levels"].(string)
			if !ok {
				log.Fatalf("error asserting hp_at_higher_levels to string")
			}
			hpAtHigherLevels = strings.ToLower(hpAtHigherLevels)

			if strings.Contains(hpAtHigherLevels, ability) {
				hpModifier = ability
			}
		}
		fmt.Println(hpModifier)

		hpAtFirstLevel := getHpAtFirstLevel(result["hp_at_first_level"])
		fmt.Println(hpAtFirstLevel)

		fmt.Println(savingThrowCharisma)
		fmt.Println(savingThrowConstitution)
		fmt.Println(savingThrowDexterity)
		fmt.Println(savingThrowIntelligence)
		fmt.Println(savingThrowStrength)
		fmt.Println(savingThrowWisdom)


	}

	return nil
}

func main() {
	dbPath := "./sql_database/mud.db"
	classesImportDbPath := "./sql_database/class_imports.db"
	SeedClasses(dbPath, classesImportDbPath)
}

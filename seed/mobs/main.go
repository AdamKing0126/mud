package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v2"
)

type MobImport struct {
	UUID        string   `yaml:"uuid"`
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Tags        []string `yaml:"tags"`
	Health      int      `yaml:"health"`
	Mana        int      `yaml:"mana"`
}

func SeedMobs() {
	db, err := sql.Open("sqlite3", "./sql_database/mud.db")
	if err != nil {
		log.Fatalf("Failed to open SQLite database: %v", err)
	} else {
		err := db.Ping()
		if err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}
		fmt.Println("Database opened successfully")
	}

	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS mobs (
			uuid VARCHAR(36) PRIMARY KEY,
			name TEXT,
			description TEXT,
			room_uuid VARCHAR(36) NULL,
			area_uuid VARCHAR(36) NULL,
			tags TEXT,
			health INTEGER,
			health_max INTEGER,
			mana INTEGER,
			mana_max INTEGER
		);
	`)

	if err != nil {
		log.Fatalf("Failed to create Areas/Rooms:  %v", err)
	}

	mobSeeds := []string{"areas/seeds/arena_mobs.yml"}
	for _, mobSeed := range mobSeeds {
		file, err := ioutil.ReadFile(mobSeed)
		if err != nil {
			log.Fatal(err)
		}

		var mobs []MobImport
		err = yaml.Unmarshal(file, &mobs)
		if err != nil {
			log.Fatal(err)
		}

		for _, mob := range mobs {
			sqlStatement := fmt.Sprintf("INSERT INTO mobs (uuid, name, description, tags, health, health_max, mana, mana_max) VALUES ('%s', '%s', '%s', '%s', '%d', '%d', '%d', '%d')", mob.UUID, mob.Name, mob.Description, mob.Tags, mob.Health, mob.Health, mob.Mana, mob.Mana)
			_, err = db.Exec(sqlStatement)
			if err != nil {
				log.Fatalf("Failed to insert mob: %v", err)
			}
		}
	}
}

func main() {
	SeedMobs()
}

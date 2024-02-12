package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	_ "github.com/mattn/go-sqlite3"
)

func CreateMobsTable(dbPath string) {
	db, err := sql.Open("sqlite3", dbPath)
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
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		area_uuid VARCHAR(36),
		room_uuid VARCHAR(36),
		alignment TEXT,
		armor_class INTEGER,
		armor_description TEXT,
		challenge_rating FLOAT,
		charisma INTEGER,
		charisma_save INTEGER,
		condition_immunities TEXT,
		constitution INTEGER,
		constitution_save INTEGER,
		damage_immunities TEXT,
		damage_resistances TEXT,
		damage_vulnerabilities TEXT,
		description TEXT,
		dexterity INTEGER,
		dexterity_save INTEGER,
		group_name TEXT,
		hp INTEGER,
		hit_dice TEXT,
		image TEXT,
		intelligence INTEGER,
		intelligence_save INTEGER,
		legendary_description TEXT,
		name TEXT,
		perception INTEGER,
		senses TEXT,
		size TEXT,
		slug TEXT,
		strength INTEGER,
		strength_save INTEGER,
		subtype TEXT,
		type TEXT,
		wisdom INTEGER,
		wisdom_save INTEGER);
	`)
	if err != nil {
		log.Fatalf("Failed to create SQLite table: %v", err)
	}

}

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

var isTest bool

func init() {
	flag.BoolVar(&isTest, "test", false, "use test db")

}

func main() {
	var dbPath string
	flag.Parse()
	if isTest {
		dbPath = "./sql_database/test_mud.db"
	} else {
		dbPath = "./sql_database/mud.db"
	}

	fmt.Printf("creating table at %s", dbPath)
	CreateMobsTable(dbPath)
	// SeedMobs()
}

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

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

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS mobs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		area_uuid VARCHAR(36),
		room_uuid VARCHAR(36),
		alignment TEXT,
		actions TEXT,
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
	defer db.Close()
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

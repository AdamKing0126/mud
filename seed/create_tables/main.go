package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func CreatePlayersTables(db *sqlx.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS players (
			uuid VARCHAR(36) PRIMARY KEY,
			character_class TEXT,
			name TEXT,
			room VARCHAR(36),
			area VARCHAR(36),
			hp INTEGER,
			movement INTEGER,
			hp_max INTEGER,
			movement_max INTEGER,
			color_profile VARCHAR(36),
			logged_in BOOLEAN DEFAULT FALSE,
			password VARCHAR(60)
		);
	`)

	if err != nil {
		log.Fatalf("Failed to create players table: %v", err)
	}
	fmt.Println("Created players table.")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS player_abilities (
			uuid VARCHAR(36) PRIMARY KEY,
			player_uuid VARCHAR(36),
			strength INTEGER,
			dexterity INTEGER,
			constitution INTEGER,
			intelligence INTEGER,
			wisdom INTEGER,
			charisma INTEGER
		);
	`)

	if err != nil {
		log.Fatalf("Failed to create player_abilities table: %v", err)
	}
	fmt.Println("Created player_abilities table.")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS player_equipments (
			uuid VARCHAR(36) PRIMARY KEY,
			player_uuid VARCHAR(36),
			Head VARCHAR(36),
			Neck VARCHAR(36),
			Chest VARCHAR(36),
			Arms VARCHAR(36),
			Hands VARCHAR(36),
			DominantHand VARCHAR(36),
			OffHand VARCHAR(36),
			Legs VARCHAR(36),
			Feet VARCHAR(36)
		);
	`)

	if err != nil {
		log.Fatalf("Failed to create player_equipments table: %v", err)
	}
	fmt.Println("Created player_equipments table.")

}

func CreateColorProfilesTable(db *sqlx.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS color_profiles (
			uuid VARCHAR(36) PRIMARY KEY,
			name TEXT,
			primary_color TEXT,
			secondary_color TEXT,
			warning_color TEXT,
			danger_color TEXT,
			title_color TEXT,
			description_color TEXT
		);
	`)

	if err != nil {
		log.Fatalf("Failed to create Color Profiles:  %v", err)
	}
	fmt.Println("Created color_profiles table.")

}

func CreateAreasAndRoomsTable(db *sqlx.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS areas (
		  uuid VARCHAR(36) PRIMARY KEY,
		  name TEXT,
		  description TEXT
		);

		CREATE TABLE IF NOT EXISTS rooms (
		  uuid VARCHAR(36) PRIMARY KEY,
		  area_uuid VARCHAR(36),
		  name TEXT,
		  description TEXT,
		  exit_north VARCHAR(36),
		  exit_south VARCHAR(36),
		  exit_east VARCHAR(36),
		  exit_west VARCHAR(36),
		  exit_up VARCHAR(36),
		  exit_down VARCHAR(36)
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create Areas/Rooms:  %v", err)
	}
	fmt.Println("Created areas table.")
	fmt.Println("Created rooms table.")
}

func CreateItemTables(db *sqlx.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS item_templates (
			uuid VARCHAR(36) PRIMARY KEY,
			name TEXT,
			description TEXT,
			equipment_slots TEXT
		);

		CREATE TABLE IF NOT EXISTS items (
			uuid VARCHAR(36) PRIMARY KEY,
			name TEXT,
			description TEXT, 
			equipment_slots TEXT
		);

		CREATE TABLE IF NOT EXISTS item_locations (
			item_uuid VARCHAR(36),
			room_uuid VARCHAR(36) NULL,
			player_uuid VARCHAR(36) NULL,
			PRIMARY KEY (item_uuid),
			FOREIGN KEY (room_uuid) REFERENCES rooms(uuid),
			FOREIGN KEY (player_uuid) REFERENCES players(uuid)
		);
	`)

	if err != nil {
		log.Fatalf("Failed to create Items:  %v", err)
	}
	fmt.Println("Created item_templates table.")
	fmt.Println("Created items table.")
	fmt.Println("Created item_locations table.")

}

func CreateMobsTable(db *sqlx.DB) {
	_, err := db.Exec(`
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
	fmt.Println("Created mobs table")
}

func CreateRacesTable(db *sqlx.DB) {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS races (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		slug TEXT,
		size TEXT,
		description TEXT,
		asi TEXT,
		subrace_name TEXT,
		subrace_slug TEXT,
		subrace_description TEXT);
	`)
	if err != nil {
		log.Fatalf("Failed to create SQLite table: %v", err)
	}
	fmt.Println("Created Races table")
}

func CreateClassesTable(db *sqlx.DB) {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS character_classes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		hit_dice TEXT,
		hp_at_first_level INTEGER,
		hp_modifier TEXT,
		name TEXT,
		saving_throw_charisma BOOL,
		saving_throw_constitution BOOL,
		saving_throw_dexterity BOOL,
		saving_throw_intelligence BOOL,
		saving_throw_strength BOOL,
		saving_throw_wisdom BOOL,
		slug TEXT,
		archetype_slug TEXT,
		archetype_name TEXT,
		archetype_description TEXT);
	`)
	if err != nil {
		log.Fatalf("Failed to create SQLite table: %v", err)
	}
	fmt.Println("Created Classes table")
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
	db, err := sqlx.Open("sqlite3", dbPath)
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

	CreateMobsTable(db)
	CreateItemTables(db)
	CreateAreasAndRoomsTable(db)
	CreateColorProfilesTable(db)
	CreatePlayersTables(db)
	CreateRacesTable(db)
	CreateClassesTable(db)
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v2"
)

type ItemImport struct {
	UUID           string   `yaml:"uuid"`
	Name           string   `yaml:"name"`
	Description    string   `yaml:"description"`
	EquipmentSlots []string `yaml:"equipment_slots"`
}

func SeedItems() {
	db, err := sqlx.Open("sqlite3", "./sql_database/mud.db")
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

	itemTemplateSeeds := []string{"areas/seeds/arena_items.yml"}
	for _, itemTemplateSeed := range itemTemplateSeeds {
		file, err := ioutil.ReadFile(itemTemplateSeed)
		if err != nil {
			log.Fatal(err)
		}

		var itemTemplates []ItemImport
		err = yaml.Unmarshal(file, &itemTemplates)
		if err != nil {
			log.Fatal(err)
		}

		for _, item := range itemTemplates {
			equipmentSlotsJSON, err := json.Marshal(item.EquipmentSlots)
			if err != nil {
				log.Fatal(err)
			}
			sqlStatement := fmt.Sprintf("INSERT INTO item_templates (uuid, name, description, equipment_slots) VALUES ('%s', '%s', '%s', '%s')", item.UUID, item.Name, item.Description, equipmentSlotsJSON)
			_, err = db.Exec(sqlStatement)
			if err != nil {
				log.Fatalf("failed to insert item: %v", err)
			}
		}
	}
}

func main() {
	SeedItems()
}

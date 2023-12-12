package main

import (
	"database/sql"
	"fmt"
	"log"
	"mud/display"

	_ "github.com/mattn/go-sqlite3"

	"github.com/google/uuid"
)

func SeedColorProfiles() {
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

	var colorProfiles = map[string]map[string]string{
		"Light Mode": {
			"primary_color":     display.BrightGreen,
			"secondary_color":   display.Green,
			"warning_color":     display.BrightYellow,
			"danger_color":      display.BrightRed,
			"title_color":       display.Reset,
			"description_color": display.Reset,
		},
		"Dark Mode": {
			"primary_color":     display.BrightCyan,
			"secondary_color":   display.Cyan,
			"warning_color":     display.BrightYellow,
			"danger_color":      display.BrightRed,
			"title_color":       display.Reset,
			"description_color": display.Reset,
		},
	}

	for name, colors := range colorProfiles {
		_, err := db.Exec(`
			INSERT INTO color_profiles (uuid, name, primary_color, secondary_color, warning_color, danger_color, title_color, description_color)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, uuid.New(), name, colors["primary_color"], colors["secondary_color"], colors["warning_color"], colors["danger_color"], colors["title_color"], colors["description_color"])

		if err != nil {
			log.Fatalf("Failed to seed Color Profiles:  %v", err)
		}
	}

}

func main() {
	SeedColorProfiles()
}

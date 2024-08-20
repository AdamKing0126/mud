package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/adamking0126/mud/display"

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

	var colorProfiles = map[string]map[string]string{
		"Light Mode": {
			"uuid":              "2c7dfd5b-d160-42e0-accb-b77d9686dbea", // this is a hack.  same as in login-logout.go
			"primary_color":     display.BrightGreen,
			"secondary_color":   display.Green,
			"warning_color":     display.BrightYellow,
			"danger_color":      display.BrightRed,
			"title_color":       display.Reset,
			"description_color": display.Reset,
		},
		"Dark Mode": {
			"uuid":              uuid.New().String(),
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
		`, colors["uuid"], name, colors["primary_color"], colors["secondary_color"], colors["warning_color"], colors["danger_color"], colors["title_color"], colors["description_color"])

		if err != nil {
			log.Fatalf("Failed to seed Color Profiles:  %v", err)
		}
	}

}

func main() {
	SeedColorProfiles()
}

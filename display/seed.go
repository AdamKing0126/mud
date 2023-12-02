package display

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
)

func SeedColorProfiles(db *sql.DB) {
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

	var colorProfiles = map[string]map[string]string{
		"Light Mode": {
			"primary_color":     BrightGreen,
			"secondary_color":   Green,
			"warning_color":     BrightYellow,
			"danger_color":      BrightRed,
			"title_color":       Reset,
			"description_color": Reset,
		},
		"Dark Mode": {
			"primary_color":     BrightCyan,
			"secondary_color":   Cyan,
			"warning_color":     BrightYellow,
			"danger_color":      BrightRed,
			"title_color":       Reset,
			"description_color": Reset,
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

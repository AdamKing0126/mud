package players

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
)

func SeedPlayers(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS players (
			uuid VARCHAR(36) PRIMARY KEY,
			name TEXT,
			room VARCHAR(36),
			area VARCHAR(36),
			health INTEGER,
			color_profile VARCHAR(36)
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create players table: %v", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM players").Scan(&count)
	if err != nil {
		log.Fatalf("Failed to query player count: %v", err)
	}

	colorProfileUUIDs := make(map[string]string)

	color_profile_rows, err := db.Query("SELECT uuid, name FROM color_profiles")
	if err != nil {
		log.Fatalf("Failed to query color profiles: %v", err)
	}
	defer color_profile_rows.Close()

	for color_profile_rows.Next() {
		var uuid, name string
		if err := color_profile_rows.Scan(&uuid, &name); err != nil {
			log.Fatalf("Failed to scan color profile: %v", err)
		}
		colorProfileUUIDs[name] = uuid
	}

	if err := color_profile_rows.Err(); err != nil {
		log.Fatalf("Failed to iterate color profiles: %v", err)
	}

	color_profile_rows.Close()

	if count == 0 {
		// Define the player data
		players := []Player{
			{
				Name:         "Reg",
				Area:         "d71e8cf1-d5ba-426c-8915-4c7f5b22e3a9",
				Room:         "189a729d-4e40-4184-a732-e2c45c66ff46",
				Health:       100,
				ColorProfile: &ColorProfile{UUID: colorProfileUUIDs["Light Mode"]},
			},
			{
				Name:         "Admin",
				Area:         "d71e8cf1-d5ba-426c-8915-4c7f5b22e3a9",
				Room:         "189a729d-4e40-4184-a732-e2c45c66ff46",
				Health:       100,
				ColorProfile: &ColorProfile{UUID: colorProfileUUIDs["Dark Mode"]},
			},
		}

		// Insert players into the database
		for _, p := range players {
			_, err := db.Exec("INSERT INTO players (uuid, name, area, room, health, color_profile) VALUES (?, ?, ?, ?, ?, ?)",
				uuid.New(), p.Name, p.Area, p.Room, p.Health, p.ColorProfile.GetUUID())
			if err != nil {
				log.Fatalf("Failed to insert player: %v", err)
			}
		}
	}
}

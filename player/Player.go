package player

import (
	"database/sql"
	"fmt"
	"log"
	"mud/areas"
	"net"

	"github.com/google/uuid"
)

type Player struct {
	UUID   string
	Name   string
	Room   string
	Area   string
	Health int
	Conn   net.Conn
}

func (player *Player) SetLocation(db *sql.DB, roomUUID string) error {
	area_rows, err := db.Query("SELECT area_uuid FROM rooms WHERE uuid=?", roomUUID)
	if err != nil {
		return fmt.Errorf("error retrieving area: %v", err)
	}
	defer area_rows.Close()

	if !area_rows.Next() {
		return fmt.Errorf("room with UUID %d does not have an area", roomUUID)
	}

	area := &areas.Area{}
	err = area_rows.Scan(&area.UUID)
	if err != nil {
		return err
	}

	player.Area = area.UUID
	player.Room = roomUUID

	stmt, err := db.Prepare("UPDATE players SET area = ?, room = ? WHERE uuid = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(area.UUID, roomUUID, player.UUID)
	if err != nil {
		return err
	}

	stmt.Close()

	return nil
}

func SeedPlayers(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS players (
			uuid VARCHAR(36) PRIMARY KEY,
			name TEXT,
			room VARCHAR(36),
			area VARCHAR(36),
			health INTEGER
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

	if count == 0 {
		// Define the player data
		players := []Player{
			{
				Name:   "Reg",
				Area:   "d71e8cf1-d5ba-426c-8915-4c7f5b22e3a9",
				Room:   "189a729d-4e40-4184-a732-e2c45c66ff46",
				Health: 100,
			},
			{
				Name:   "Admin",
				Area:   "d71e8cf1-d5ba-426c-8915-4c7f5b22e3a9",
				Room:   "189a729d-4e40-4184-a732-e2c45c66ff46",
				Health: 100,
			},
		}

		// Insert players into the database
		for _, p := range players {
			_, err := db.Exec("INSERT INTO players (uuid, name, area, room, health) VALUES (?, ?, ?, ?, ?)",
				uuid.New(), p.Name, p.Area, p.Room, p.Health)
			if err != nil {
				log.Fatalf("Failed to insert player: %v", err)
			}
		}
	}
}

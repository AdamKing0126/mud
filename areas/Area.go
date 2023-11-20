package areas

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Area struct {
	UUID        string
	Name        string
	Description string
	Rooms       []Room
}

type AreaInfo struct {
	UUID        string
	Name        string
	Description string
}

type AreaImport struct {
	UUID        string       `yaml:"uuid"`
	Name        string       `yaml:"name"`
	Description string       `yaml:"description"`
	Rooms       []RoomImport `yaml:"rooms"`
}

func SeedAreasAndRooms(db *sql.DB) {
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

	areaSeeds := []string{"areas/arena.yml", "areas/street.yml", "areas/glade.yml"}

	for _, areaSeed := range areaSeeds {
		// Read the YAML file.
		file, err := ioutil.ReadFile(areaSeed)
		if err != nil {
			log.Fatal(err)
		}

		// Parse the YAML file.
		var area AreaImport
		err = yaml.Unmarshal(file, &area)
		if err != nil {
			log.Fatal(err)
		}
		var count int

		err = db.QueryRow("SELECT COUNT(*) FROM areas WHERE uuid=?", area.UUID).Scan(&count)
		if err != nil {
			log.Fatalf("Failed to query area count: %v", err)
		}

		if count == 0 {
			_, err := db.Exec("INSERT INTO areas (uuid, name, description) VALUES (?, ?, ?)",
				area.UUID, area.Name, area.Description)
			if err != nil {
				log.Fatalf("Failed to insert area: %v", err)
			}

			for _, room := range area.Rooms {
				sqlStatement := fmt.Sprintf("INSERT INTO rooms (uuid, area_uuid, name, description, exit_north, exit_south, exit_west, exit_east, exit_up, exit_down) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')", room.UUID, area.UUID, room.Name, room.Description, room.Exits["north"], room.Exits["south"], room.Exits["west"], room.Exits["east"], room.Exits["up"], room.Exits["down"])
				_, err := db.Exec(sqlStatement)
				if err != nil {
					log.Fatalf("Failed to insert room: %v", err)
				}
			}

		}
	}
}

func LoadAreaFromDB(db *sql.DB, areaUUID string) (*Area, error) {
	// Query the database for the area data.
	area_rows, err := db.Query("SELECT uuid, name, description FROM areas where uuid=?", areaUUID)
	if err != nil {
		return nil, err
	}

	// Make sure that the area exists.
	if !area_rows.Next() {
		return nil, fmt.Errorf("Area with UUID %d does not exist", areaUUID)
	}

	// Create a new Area struct.
	area := &Area{
		UUID: areaUUID,
	}

	// Scan the row and populate the Area struct.
	err = area_rows.Scan(&area.UUID, &area.Name, &area.Description)
	if err != nil {
		return nil, err
	}

	// Close the rows.
	if err := area_rows.Close(); err != nil {
		return nil, err
	}

	var rooms []Room
	room_rows, err := db.Query("SELECT uuid, area_uuid, name, description, exit_north, exit_south, exit_west, exit_east, exit_up, exit_down FROM rooms where area_uuid=?", areaUUID)
	if err != nil {
		return nil, err
	}

	for room_rows.Next() {
		var room Room
		err := room_rows.Scan(&room.UUID, &room.AreaUUID, &room.Name, &room.Description, &room.Exits.North, &room.Exits.South, &room.Exits.West, &room.Exits.East, &room.Exits.Up, &room.Exits.Down)
		if err != nil {
			return nil, err
		}

		rooms = append(rooms, room)

	}
	area.Rooms = rooms

	if err := room_rows.Close(); err != nil {
		return nil, err
	}

	return area, nil
}

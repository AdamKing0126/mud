package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"mud/items"
	"mud/mobs"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

type RoomImport struct {
	UUID        string            `yaml:"uuid"`
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Exits       map[string]string `yaml:"exits"`
	Mobs        []string          `yaml:"mobs"`
	Items       []string          `yaml:"items"`
}

type AreaImport struct {
	UUID        string       `yaml:"uuid"`
	Name        string       `yaml:"name"`
	Description string       `yaml:"description"`
	Rooms       []RoomImport `yaml:"rooms"`
}

func SeedAreasAndRooms(dbPath string, monstersImportDbPath string) error {
	// db, err := sql.Open("sqlite3", dbPath)
	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open SQLite database: %v", err)
	} else {
		err := db.Ping()
		if err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}
		fmt.Println("Mud Database opened successfully")
	}
	defer db.Close()

	monstersDB, err := sqlx.Connect("sqlite3", monstersImportDbPath)
	if err != nil {
		log.Fatalf("Failed to open Monster Imports database: %v", err)
	} else {
		err := db.Ping()
		if err != nil {
			log.Fatalf("Failed to ping database Monster Imports: %v", err)
		}
		fmt.Println("Monster Imports Database opened successfully")
	}
	defer monstersDB.Close()

	_, err = db.Exec(`
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
		log.Fatalf("Failed to create Areas/Rooms:  %v", err)
	}

	areaSeeds := []string{"areas/seeds/arena.yml", "areas/seeds/street.yml", "areas/seeds/glade.yml"}

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

			// for idx, room := range area.Rooms {
			for _, room := range area.Rooms {
				sqlStatement := fmt.Sprintf("INSERT INTO rooms (uuid, area_uuid, name, description, exit_north, exit_south, exit_west, exit_east, exit_up, exit_down) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')", room.UUID, area.UUID, room.Name, room.Description, room.Exits["north"], room.Exits["south"], room.Exits["west"], room.Exits["east"], room.Exits["up"], room.Exits["down"])
				_, err := db.Exec(sqlStatement)
				if err != nil {
					log.Fatalf("Failed to insert room: %v", err)
				}
				for _, slug := range room.Mobs {
					// get the mob by slug from the monsters_import database

					var tableNames []string
					err = monstersDB.Select(&tableNames, "SELECT name from sqlite_master WHERE type='table';")
					if err != nil {
						log.Fatalf("%v", err)
					}

					row := monstersDB.QueryRowx("SELECT * from mob_imports where slug = ?", slug)

					// have to create a result interface, in order to ignore columns in the
					// table which are not present on the struct
					result := make(map[string]interface{})
					err = row.MapScan(result)
					if err != nil {
						log.Fatal(err)
					}

					var mob mobs.Mob
					err = mapstructure.Decode(result, &mob)

					if err != nil {
						log.Fatalf("failed to fetch mob from mob_imports: %v", err)
					}

					mob.RoomUUID = room.UUID
					mob.AreaUUID = area.UUID

					fmt.Println("yo")
					mapper := reflectx.NewMapperFunc("db", strings.ToLower)

					fieldInfos := mapper.TypeMap(reflect.TypeOf(mobs.Mob{})).Names
					fields := make([]string, 0, len(fieldInfos))
					for field := range fieldInfos {
						fields = append(fields, field)
					}

					columns := strings.Join(fields, ", :")
					columns = ":" + columns

					query := fmt.Sprintf("INSERT INTO mobs (%s) VALUES (%s)", strings.Join(fields, ", "), columns)

					_, err := db.NamedExec(query, mob)
					if err != nil {
						log.Fatalf("failed to insert mob into mobs: %v", err)
					}

				}
				for idx, item := range room.Items {
					newItem, err := items.NewItemFromTemplate(db, item)
					if err != nil {
						log.Fatalf("error creating item from template, %v", err)
					}

					_, item_location_err := db.Exec("INSERT INTO item_locations (item_uuid, room_uuid, player_uuid) VALUES (?, ?, NULL)", newItem.UUID, area.Rooms[idx].UUID)
					if item_location_err != nil {
						log.Fatalf("Failed to insert item location: %v", item_location_err)
					}

				}
			}

			item_uuid := uuid.New().String()
			_, item_err := db.Exec("INSERT INTO items (uuid, name, description, equipment_slots) VALUES (?, ?, ?, ?)", item_uuid, "sword", "A sword", "DominantHand,OffHand")
			if item_err != nil {
				log.Fatalf("Failed to insert item: %v", item_err)
			}

			_, item_location_err := db.Exec("INSERT INTO item_locations (item_uuid, room_uuid, player_uuid) VALUES (?, ?, NULL)", item_uuid, area.Rooms[0].UUID)
			if item_location_err != nil {
				log.Fatalf("Failed to insert item location: %v", item_location_err)
			}
		}
	}
	return nil
}

func main() {
	dbPath := "./sql_database/mud.db"
	monstersImportDbPath := "./sql_database/monster_imports.db"
	SeedAreasAndRooms(dbPath, monstersImportDbPath)
}

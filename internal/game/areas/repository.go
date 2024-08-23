package areas

import (
	"context"
	"fmt"

	"github.com/adamking0126/mud/internal/game/items"
	"github.com/adamking0126/mud/pkg/database"
)

type Repository struct {
	db database.DB
}

type AreaData struct {
	RoomUUID        string
	AreaUUID        string
	AreaName        string
	AreaDescription string
}

func NewRepository(db database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetRoomFromDBAndLoadArea(ctx context.Context, area *Area, targetRoomUUID string) *Room {
	var retrievedRoom *Room
	queryString := `
	SELECT r.UUID, r.area_uuid, r.name, r.description,
		r.exit_north, r.exit_south, r.exit_east, r.exit_west,
		r.exit_up, r.exit_down,
		a.UUID AS area_uuid, a.name AS area_name, a.description AS area_description
	FROM rooms r
	LEFT JOIN areas a ON r.area_uuid = a.UUID
	WHERE a.UUID = ?`
	rows, err := r.db.Query(ctx, queryString, area.UUID)
	if err != nil {
		fmt.Printf("Error querying rows: %v", err)
	}
	for rows.Next() {
		var roomUUID, roomAreaUUID, roomName, roomDescription, exitNorth, exitSouth, exitEast, exitWest, exitUp, exitDown, areaUUID, areaName, areaDescription string
		err := rows.Scan(&roomUUID, &roomAreaUUID, &roomName, &roomDescription, &exitNorth, &exitSouth, &exitEast, &exitWest, &exitUp, &exitDown, &areaUUID, &areaName, &areaDescription)
		if err != nil {
			fmt.Printf("error scanning rows: %v", err)
		}
		room := NewRoomWithAreaInfo(roomUUID, roomAreaUUID, roomName, roomDescription, area.Name, area.Description, exitNorth, exitSouth, exitEast, exitWest, exitUp, exitDown)
		area.Rooms = (append(area.Rooms, room))
		if roomUUID == targetRoomUUID {
			retrievedRoom = room
		}
	}
	defer rows.Close()

	return retrievedRoom
}

func (r *Repository) GetRoomByUUIDFromDB(ctx context.Context, roomUUID string) *Room {
	query := `
			SELECT r.UUID, r.area_uuid, r.name, r.description,
				r.exit_north, r.exit_south, r.exit_east, r.exit_west,
				r.exit_up, r.exit_down,
				a.UUID AS area_uuid, a.name AS area_name, a.description AS area_description
			FROM rooms r
			LEFT JOIN areas a ON r.area_uuid = a.UUID
			WHERE r.UUID = ?`

	room_rows, err := r.db.Query(ctx, query, roomUUID)
	if err != nil {
		return nil
	}

	defer room_rows.Close()
	if !room_rows.Next() {
		return nil
	}

	var northExitUUID, southExitUUID, eastExitUUID, westExitUUID, upExitUUID, downExitUUID string
	room := &Room{Exits: &ExitInfo{}}
	err = room_rows.Scan(
		&room.UUID, &room.AreaUUID, &room.Name, &room.Description,
		&northExitUUID, &southExitUUID, &eastExitUUID, &westExitUUID,
		&upExitUUID, &downExitUUID,
		&room.Area.UUID, &room.Area.Name, &room.Area.Description,
	)
	if err != nil {
		return nil
	}

	exitUUIDs := map[string]*string{
		"North": &northExitUUID,
		"South": &southExitUUID,
		"West":  &westExitUUID,
		"East":  &eastExitUUID,
		"Down":  &downExitUUID,
		"Up":    &upExitUUID,
	}

	for direction, uuid := range exitUUIDs {
		if *uuid != "" {
			switch direction {
			case "North":
				room.Exits.North = &Room{UUID: *uuid}
			case "South":
				room.Exits.South = &Room{UUID: *uuid}
			case "West":
				room.Exits.West = &Room{UUID: *uuid}
			case "East":
				room.Exits.East = &Room{UUID: *uuid}
			case "Down":
				room.Exits.Down = &Room{UUID: *uuid}
			case "Up":
				room.Exits.Up = &Room{UUID: *uuid}
			}
		}
	}

	return room
}

// relaces room.AddItem() - need to find where this is called
func (r *Repository) AddItemToRoom(ctx context.Context, room *Room, item *items.Item) error {
	query := fmt.Sprintf("UPDATE item_locations SET room_uuid = '%s', player_uuid = '' WHERE item_uuid = '%s'", room.UUID, item.UUID)
	err := r.db.Exec(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetAreasFromDB(ctx context.Context) ([]AreaData, error) {
	var areaData []AreaData

	queryString := `
			SELECT r.uuid, a.uuid, a.name, a.description 
			FROM rooms r
			JOIN areas a ON r.area_uuid = a.uuid;
	`
	rows, err := r.db.Query(ctx, queryString)
	if err != nil {
		return nil, fmt.Errorf("error retrieving areas/rooms: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var data AreaData
		err := rows.Scan(&data.RoomUUID, &data.AreaUUID, &data.AreaName, &data.AreaDescription)
		if err != nil {
			return nil, fmt.Errorf("error scanning areas/rooms: %v", err)
		}
		areaData = append(areaData, data)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over areas/rooms: %v", err)
	}

	return areaData, nil
}

func (r *Repository) GetRoomCount(ctx context.Context, areaUUID string) (int, error) {
	query := `SELECT COUNT(*) FROM rooms WHERE area_uuid = ?`
	var count int
	err := r.db.QueryRow(ctx, query, areaUUID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error retrieving room count: %v", err)
	}
	return count, nil
}

func (r *Repository) GetRoomAndLoadAreaWithRooms(ctx context.Context, area *Area, targetRoomUUID string, followExits bool) *Room {
	queryString := `
	SELECT r.UUID, r.area_uuid, r.name, r.description,
		r.exit_north, r.exit_south, r.exit_east, r.exit_west,
		r.exit_up, r.exit_down,
		a.UUID AS area_uuid, a.name AS area_name, a.description AS area_description
	FROM rooms r
	LEFT JOIN areas a ON r.area_uuid = a.UUID
	WHERE a.UUID = ?`
	rows, err := r.db.Query(ctx, queryString, area.UUID)
	if err != nil {
		fmt.Printf("Error querying rows: %v", err)
	}
	for rows.Next() {
		var roomUUID, roomAreaUUID, roomName, roomDescription, exitNorth, exitSouth, exitEast, exitWest, exitUp, exitDown, areaUUID, areaName, areaDescription string
		err := rows.Scan(&roomUUID, &roomAreaUUID, &roomName, &roomDescription, &exitNorth, &exitSouth, &exitEast, &exitWest, &exitUp, &exitDown, &areaUUID, &areaName, &areaDescription)
		if err != nil {
			fmt.Printf("error scanning rows: %v", err)
		}
		room := NewRoomWithAreaInfo(roomUUID, roomAreaUUID, roomName, roomDescription, area.Name, area.Description, exitNorth, exitSouth, exitEast, exitWest, exitUp, exitDown)
		area.Rooms = (append(area.Rooms, room))
		if roomUUID == targetRoomUUID {
			return room
		}
	}
	defer rows.Close()
	return nil
}

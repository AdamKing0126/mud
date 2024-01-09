package areas

import (
	"database/sql"
	"fmt"
)

func GetRoomFromDB(roomUUID string, db *sql.DB) (*Room, error) {
	query := `
		SELECT r.UUID, r.area_uuid, r.name, r.description,
			r.exit_north, r.exit_south, r.exit_east, r.exit_west,
			r.exit_up, r.exit_down,
			a.UUID AS area_uuid, a.name AS area_name, a.description AS area_description
		FROM rooms r
		LEFT JOIN areas a ON r.area_uuid = a.UUID
		WHERE r.UUID = ?`

	room_rows, err := db.Query(query, roomUUID)
	if err != nil {
		return nil, err
	}

	defer room_rows.Close()
	if !room_rows.Next() {
		return nil, fmt.Errorf("room with UUID %s does not exist", roomUUID)
	}

	var northExitUUID, southExitUUID, eastExitUUID, westExitUUID, upExitUUID, downExitUUID string
	room := &Room{Exits: ExitInfo{}}
	err = room_rows.Scan(
		&room.UUID, &room.AreaUUID, &room.Name, &room.Description,
		&northExitUUID, &southExitUUID, &eastExitUUID, &westExitUUID,
		&upExitUUID, &downExitUUID,
		&room.Area.UUID, &room.Area.Name, &room.Area.Description,
	)
	if err != nil {
		return nil, err
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

	return room, nil
}

package areas

import (
	"database/sql"
	"fmt"
	"mud/interfaces"
)

type Area struct {
	UUID        string
	Name        string
	Description string
	Rooms       []Room
	Channel     chan Action
}

func (a *Area) GetUUID() string {
	return a.UUID
}

func (a *Area) GetName() string {
	return a.Name
}

func (a *Area) GetDescription() string {
	return a.Description
}

type AreaInfo struct {
	UUID        string
	Name        string
	Description string
}

func NewArea() interfaces.AreaInterface {
	return &Area{}
}

// TODO: No usages of this function.  Remove?
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

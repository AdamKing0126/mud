package world_state

import (
	"fmt"
	"mud/areas"
	"mud/players"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type WorldState struct {
	Areas         map[string]*areas.Area
	RoomToAreaMap map[string]string
	DB            *sqlx.DB
}

func NewWorldState(areas map[string]*areas.Area, roomToAreaMap map[string]string, db *sqlx.DB) *WorldState {
	return &WorldState{Areas: areas, RoomToAreaMap: roomToAreaMap, DB: db}
}

func (worldState *WorldState) RemovePlayerFromRoom(roomUUID string, player *players.Player) error {
	areaUUID := worldState.RoomToAreaMap[roomUUID]
	area := worldState.Areas[areaUUID]
	room, err := area.GetRoomByUUID(roomUUID)
	if err != nil {
		return err
	}

	err = room.RemovePlayer(player)
	if err != nil {
		return err
	}

	return nil
}

func (worldState *WorldState) AddPlayerToRoom(roomUUID string, player *players.Player) error {
	areaUUID := worldState.RoomToAreaMap[roomUUID]
	area := worldState.Areas[areaUUID]
	room, err := area.GetRoomByUUID(roomUUID)
	if err != nil {
		return err
	}
	room.AddPlayer(player)
	player.RoomUUID = roomUUID

	return nil
}

func (worldState *WorldState) GetRoom(roomUUID string, followExits bool) *areas.Room {
	areaUUID := worldState.RoomToAreaMap[roomUUID]

	queryString := `
		SELECT COUNT(*) AS room_count
		FROM rooms
		WHERE area_uuid = ?`
	rows, err := worldState.DB.Query(queryString, areaUUID)
	if err != nil {
		fmt.Printf("Error querying rows: %v", err)
	}
	defer rows.Close()

	var numberOfRoomsInArea int
	for rows.Next() {
		err := rows.Scan(&numberOfRoomsInArea)
		if err != nil {
			fmt.Printf("Error scanning rows: %v", err)
		}
	}

	area := worldState.Areas[areaUUID]
	if len(area.Rooms) == numberOfRoomsInArea {
		for idx := range area.Rooms {
			if area.Rooms[idx].UUID == roomUUID {
				return area.Rooms[idx]
			}
		}
	}
	return retrieveRoomFromDB(worldState.DB, area, roomUUID, followExits)
}

func (worldState *WorldState) GetArea(areaUUID string) *areas.Area {
	return worldState.Areas[areaUUID]
}

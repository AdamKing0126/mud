package world_state

import (
	"database/sql"
	"fmt"
	"mud/interfaces"
)

type WorldState struct {
	Areas         map[string]interfaces.Area
	RoomToAreaMap map[string]string
	DB            *sql.DB
}

func NewWorldState(areas map[string]interfaces.Area, roomToAreaMap map[string]string, db *sql.DB) *WorldState {
	return &WorldState{Areas: areas, RoomToAreaMap: roomToAreaMap, DB: db}
}

func (worldState *WorldState) RemovePlayerFromRoom(roomUUID string, player interfaces.Player) error {
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

func (worldState *WorldState) AddPlayerToRoom(roomUUID string, player interfaces.Player) error {
	areaUUID := worldState.RoomToAreaMap[roomUUID]
	area := worldState.Areas[areaUUID]
	room, err := area.GetRoomByUUID(roomUUID)
	if err != nil {
		return err
	}
	room.AddPlayer(player)
	player.SetRoom(room)

	return nil
}

func (worldState *WorldState) GetRoom(roomUUID string, followExits bool) interfaces.Room {
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
	rooms := area.GetRooms()
	if len(rooms) == numberOfRoomsInArea {
		for idx := range rooms {
			if rooms[idx].GetUUID() == roomUUID {
				return rooms[idx]
			}
		}
	}
	return retrieveRoomFromDB(worldState.DB, area, roomUUID, followExits)
}

func (worldState *WorldState) GetArea(areaUUID string) interfaces.Area {
	return worldState.Areas[areaUUID]
}

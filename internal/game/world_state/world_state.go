package world_state

import (
	"context"
	"fmt"

	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/internal/game/players"
	"github.com/adamking0126/mud/pkg/database"

	_ "github.com/mattn/go-sqlite3"
)

type WorldState struct {
	Areas         map[string]*areas.Area
	RoomToAreaMap map[string]string
	DB            database.DB
	AreasService  *areas.Service
}

func NewWorldState(ctx context.Context, areas map[string]*areas.Area, roomToAreaMap map[string]string, db database.DB, areasService *areas.Service) *WorldState {
	return &WorldState{Areas: areas, RoomToAreaMap: roomToAreaMap, DB: db, AreasService: areasService}
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

func (worldState *WorldState) GetRoom(ctx context.Context, roomUUID string, followExits bool) *areas.Room {
	areaUUID := worldState.RoomToAreaMap[roomUUID]

	queryString := `
		SELECT COUNT(*) AS room_count
		FROM rooms
		WHERE area_uuid = ?`
	rows, err := worldState.DB.Query(ctx, queryString, areaUUID)
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
	return retrieveRoomFromDB(ctx, worldState.DB, area, roomUUID, followExits)
}

func (worldState *WorldState) GetArea(ctx context.Context, areaUUID string) *areas.Area {
	return worldState.Areas[areaUUID]
}

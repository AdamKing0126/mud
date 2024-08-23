package world_state

import (
	"context"

	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/internal/game/players"
	"github.com/adamking0126/mud/pkg/database"

	_ "github.com/mattn/go-sqlite3"
)

// WorldState is, shockingly, the state of the game world.
// Chiefly, we maintain a map of area UUIDs to areas.
// All changes that happen in each room of the MUD are reflected in the WorldState.
// We also maintain a map of room UUIDs to area UUIDs.
// This allows us to look up the area of a room, which is useful for things like
// moving players between rooms.
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

func (worldState *WorldState) GetArea(ctx context.Context, areaUUID string) *areas.Area {
	return worldState.Areas[areaUUID]
}

package world_state

import (
	"context"
	"fmt"

	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/internal/game/items"
	"github.com/adamking0126/mud/internal/game/mobs"
	"github.com/adamking0126/mud/internal/game/players"
	"github.com/adamking0126/mud/pkg/database"
)

func getRoomFromAreaByUUID(ctx context.Context, area *areas.Area, roomUUID string) *areas.Room {
	for _, room := range area.Rooms {
		if room.UUID == roomUUID {
			return room
		}
	}
	fmt.Printf("Room %s not found for area %s - %s\n", roomUUID, area.Name, area.UUID)
	return nil
}

func retrieveRoomFromDB(ctx context.Context, db database.DB, area *areas.Area, requestedRoomUUID string, followExits bool) *areas.Room {
	var retrievedRoom *areas.Room
	queryString := `
	SELECT r.UUID, r.area_uuid, r.name, r.description,
		r.exit_north, r.exit_south, r.exit_east, r.exit_west,
		r.exit_up, r.exit_down,
		a.UUID AS area_uuid, a.name AS area_name, a.description AS area_description
	FROM rooms r
	LEFT JOIN areas a ON r.area_uuid = a.UUID
	WHERE a.UUID = ?`
	rows, err := db.Query(ctx, queryString, area.UUID)
	if err != nil {
		fmt.Printf("Error querying rows: %v", err)
	}
	for rows.Next() {
		var roomUUID, roomAreaUUID, roomName, roomDescription, exitNorth, exitSouth, exitEast, exitWest, exitUp, exitDown, areaUUID, areaName, areaDescription string
		err := rows.Scan(&roomUUID, &roomAreaUUID, &roomName, &roomDescription, &exitNorth, &exitSouth, &exitEast, &exitWest, &exitUp, &exitDown, &areaUUID, &areaName, &areaDescription)
		if err != nil {
			fmt.Printf("error scanning rows: %v", err)
		}
		room := areas.NewRoomWithAreaInfo(roomUUID, roomAreaUUID, roomName, roomDescription, area.Name, area.Description, exitNorth, exitSouth, exitEast, exitWest, exitUp, exitDown)
		area.Rooms = (append(area.Rooms, room))
		if roomUUID == requestedRoomUUID {
			retrievedRoom = room
		}
	}
	defer rows.Close()

	// now that we have loaded all the rooms in the area, we can go back and hook up all the
	// exits.  if one of the exits happens to exist in a different area, we can make a db query to retrieve that one.
	for _, roomInArea := range area.Rooms {
		if followExits {
			setExits(ctx, db, area, roomInArea)
		}
		setItems(ctx, db, roomInArea)
		setPlayers(ctx, db, roomInArea)
		setMobs(ctx, db, roomInArea)
	}
	return retrievedRoom
}

func getExitRoom(ctx context.Context, area *areas.Area, room *areas.Room, db database.DB) *areas.Room {
	if room != nil {
		if room.Name == "" {
			exitRoom := getRoomFromAreaByUUID(ctx, area, room.UUID)
			// TODO Adam - looks like this is not triggering, because we end up loading the adjoining areas
			// from the db.  but not the adjoining-adjoining areas, if that makes sense.
			if exitRoom == nil {
				rm, err := areas.GetRoomFromDB(ctx, db, room.UUID)
				if err != nil {
					fmt.Printf("error getting room: %v", err)
				}
				return rm
			}
			return exitRoom
		}
	}
	return nil
}

func setExits(ctx context.Context, db database.DB, area *areas.Area, roomInArea *areas.Room) {
	exits := roomInArea.Exits
	exitInfo := areas.ExitInfo{}
	exitInfo.South = getExitRoom(ctx, area, exits.South, db)
	exitInfo.North = getExitRoom(ctx, area, exits.North, db)
	exitInfo.West = getExitRoom(ctx, area, exits.West, db)
	exitInfo.East = getExitRoom(ctx, area, exits.East, db)
	exitInfo.Up = getExitRoom(ctx, area, exits.Up, db)
	exitInfo.Down = getExitRoom(ctx, area, exits.Down, db)
	// TODO wtf is this even doing?
	// am I just refreshing the exits from the db?
	// roomInArea.SetExits(exitInfo)
	roomInArea.Exits = &exitInfo
}

func setMobs(ctx context.Context, db database.DB, roomInArea *areas.Room) {
	// retrieve the mobs and attach them to the room in the WorldState
	mobsInRoom, err := mobs.GetMobsInRoom(ctx, db, roomInArea.UUID)
	if err != nil {
		fmt.Printf("error retrieving mobs: %v", err)
	}
	roomInArea.Mobs = mobsInRoom
}

func setPlayers(ctx context.Context, db database.DB, roomInArea *areas.Room) {
	// retrieve the players and attach them to the room in the WorldState
	playersInRoom, err := players.GetPlayersInRoom(ctx, db, roomInArea.UUID)
	if err != nil {
		fmt.Printf("error retrieving players: %v", err)
	}
	roomInArea.Players = playersInRoom
}

func setItems(ctx context.Context, db database.DB, roomInArea *areas.Room) {
	// retrieve the items and attach them to the room in the WorldState
	itemsInRoom, err := items.GetItemsInRoom(ctx, db, roomInArea.UUID)
	if err != nil {
		fmt.Printf("error retrieving items: %v", err)
	}
	roomInArea.Items = itemsInRoom
}

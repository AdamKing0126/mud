package world_state

import (
	"fmt"
	"mud/areas"
	"mud/items"
	"mud/mobs"
	"mud/players"

	"github.com/jmoiron/sqlx"
)

func getRoomFromAreaByUUID(area areas.Area, roomUUID string) *areas.Room {
	for _, room := range area.GetRooms() {
		if room.GetUUID() == roomUUID {
			return &room
		}
	}
	fmt.Printf("Room %s not found for area %s - %s\n", roomUUID, area.GetName(), area.GetUUID())
	return nil
}

func retrieveRoomFromDB(db *sqlx.DB, area areas.Area, requestedRoomUUID string, followExits bool) areas.Room {
	var retrievedRoom areas.Room
	queryString := `
	SELECT r.UUID, r.area_uuid, r.name, r.description,
		r.exit_north, r.exit_south, r.exit_east, r.exit_west,
		r.exit_up, r.exit_down,
		a.UUID AS area_uuid, a.name AS area_name, a.description AS area_description
	FROM rooms r
	LEFT JOIN areas a ON r.area_uuid = a.UUID
	WHERE a.UUID = ?`
	rows, err := db.Query(queryString, area.GetUUID())
	if err != nil {
		fmt.Printf("Error querying rows: %v", err)
	}
	for rows.Next() {
		var roomUUID, roomAreaUUID, roomName, roomDescription, exitNorth, exitSouth, exitEast, exitWest, exitUp, exitDown, areaUUID, areaName, areaDescription string
		err := rows.Scan(&roomUUID, &roomAreaUUID, &roomName, &roomDescription, &exitNorth, &exitSouth, &exitEast, &exitWest, &exitUp, &exitDown, &areaUUID, &areaName, &areaDescription)
		if err != nil {
			fmt.Printf("error scanning rows: %v", err)
		}
		room := areas.NewRoomWithAreaInfo(roomUUID, roomAreaUUID, roomName, roomDescription, area.GetName(), area.GetDescription(), exitNorth, exitSouth, exitEast, exitWest, exitUp, exitDown)
		area.SetRooms(append(area.GetRooms(), room))
		if roomUUID == requestedRoomUUID {
			retrievedRoom = room
		}
	}
	defer rows.Close()

	// now that we have loaded all the rooms in the area, we can go back and hook up all the
	// exits.  if one of the exits happens to exist in a different area, we can make a db query to retrieve that one.
	for _, roomInArea := range area.GetRooms() {
		if followExits {
			setExits(db, area, roomInArea)
		}
		setItems(db, &roomInArea)
		setPlayers(db, &roomInArea)
		setMobs(db, &roomInArea)
	}
	return retrievedRoom
}

func getExitRoom(area areas.Area, room *areas.Room, db *sqlx.DB) *areas.Room {
	if room != nil {
		if room.GetName() == "" {
			exitRoom := getRoomFromAreaByUUID(area, room.GetUUID())
			// TODO Adam - looks like this is not triggering, because we end up loading the adjoining areas
			// from the db.  but not the adjoining-adjoining areas, if that makes sense.
			if exitRoom == nil {
				rm, err := areas.GetRoomFromDB(room.GetUUID(), db)
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

func setExits(db *sqlx.DB, area areas.Area, roomInArea areas.Room) {
	exits := roomInArea.GetExits()
	exitInfo := areas.ExitInfo{}
	exitInfo.South = getExitRoom(area, exits.GetSouth(), db)
	exitInfo.North = getExitRoom(area, exits.GetNorth(), db)
	exitInfo.West = getExitRoom(area, exits.GetWest(), db)
	exitInfo.East = getExitRoom(area, exits.GetEast(), db)
	exitInfo.Up = getExitRoom(area, exits.GetUp(), db)
	exitInfo.Down = getExitRoom(area, exits.GetDown(), db)
	roomInArea.SetExits(exitInfo)
}

func setMobs(db *sqlx.DB, roomInArea *areas.Room) {
	// retrieve the mobs and attach them to the room in the WorldState
	mobsInRoom, err := mobs.GetMobsInRoom(db, roomInArea.GetUUID())
	if err != nil {
		fmt.Printf("error retrieving mobs: %v", err)
	}
	roomInArea.SetMobs(mobsInRoom)
}

func setPlayers(db *sqlx.DB, roomInArea *areas.Room) {
	// retrieve the players and attach them to the room in the WorldState
	playersInRoom, err := players.GetPlayersInRoom(db, roomInArea.GetUUID())
	if err != nil {
		fmt.Printf("error retrieving players: %v", err)
	}
	roomInArea.SetPlayers(playersInRoom)
}

func setItems(db *sqlx.DB, roomInArea *areas.Room) {
	// retrieve the items and attach them to the room in the WorldState
	itemsInRoom, err := items.GetItemsInRoom(db, roomInArea.GetUUID())
	if err != nil {
		fmt.Printf("error retrieving items: %v", err)
	}
	roomInArea.SetItems(itemsInRoom)
}

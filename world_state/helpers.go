package world_state

import (
	"database/sql"
	"fmt"
	"mud/areas"
	"mud/interfaces"
	"mud/items"
	"mud/players"
)

func getRoomFromAreaByUUID(area interfaces.Area, roomUUID string) interfaces.Room {
	for _, room := range area.GetRooms() {
		if room.GetUUID() == roomUUID {
			return room
		}
	}
	fmt.Println("Room not found")
	return nil
}

func retrieveRoomFromDB(db *sql.DB, area interfaces.Area, roomUUID string, followExits bool) interfaces.Room {
	var retrievedRoom interfaces.Room
	fmt.Println("rooms in memory does not match rooms in DB, must make a call to the DB")
	queryString := `
            SELECT r.UUID, r.area_uuid, r.name, r.description,
                r.exit_north, r.exit_south, r.exit_east, r.exit_west,
                r.exit_up, r.exit_down,
                a.UUID AS area_uuid, a.name AS area_name, a.description AS area_description
            FROM rooms r
            LEFT JOIN areas a ON r.area_uuid = a.UUID
            WHERE a.UUID = ? AND r.UUID = ?`
	row := db.QueryRow(queryString, area.GetUUID(), roomUUID)
	var uuid, roomAreaUUID, roomName, roomDescription, exitNorth, exitSouth, exitEast, exitWest, exitUp, exitDown, areaUUID, areaName, areaDescription string
	err := row.Scan(&uuid, &roomAreaUUID, &roomName, &roomDescription, &exitNorth, &exitSouth, &exitEast, &exitWest, &exitUp, &exitDown, &areaUUID, &areaName, &areaDescription)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle no result
			fmt.Println("No rows were returned!")
			return nil
		} else {
			// Handle other errors
			fmt.Printf("error scanning row: %v", err)
			return nil
		}
	}
	room := areas.NewRoomWithAreaInfo(uuid, roomAreaUUID, roomName, roomDescription, area.GetName(), area.GetDescription(), exitNorth, exitSouth, exitEast, exitWest, exitUp, exitDown)
	area.SetRooms(append(area.GetRooms(), room))
	if room.UUID == roomUUID {
		retrievedRoom = room
	}

	// now that we have loaded all the rooms in the area, we can go back and hook up all the
	// exits.  if one of the exits happens to exist in a different area, we can make a db query to retrieve that one.
	for idx, roomInArea := range area.GetRooms() {
		if followExits {
			setExits(db, area, roomInArea)
		}
		setItems(db, roomInArea)
		setPlayers(db, roomInArea)

		area.SetRoomAtIndex(idx, roomInArea)
	}
	return retrievedRoom
}

func getExitRoom(area interfaces.Area, room interfaces.Room, db *sql.DB) interfaces.Room {
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

func setExits(db *sql.DB, area interfaces.Area, roomInArea interfaces.Room) {
	exits := roomInArea.GetExits()
	exitInfo := areas.ExitInfo{}
	exitInfo.South = getExitRoom(area, exits.GetSouth(), db)
	exitInfo.North = getExitRoom(area, exits.GetNorth(), db)
	exitInfo.West = getExitRoom(area, exits.GetWest(), db)
	exitInfo.East = getExitRoom(area, exits.GetEast(), db)
	exitInfo.Up = getExitRoom(area, exits.GetUp(), db)
	exitInfo.Down = getExitRoom(area, exits.GetDown(), db)
}

func setPlayers(db *sql.DB, roomInArea interfaces.Room) {
	// retrieve the players and attach them to the room in the WorldState
	playersInRoom, err := players.GetPlayersInRoom(db, roomInArea.GetUUID())
	if err != nil {
		fmt.Printf("error retrieving players: %v", err)
	}
	playerInterfaces := make([]interfaces.Player, len(playersInRoom))
	for i, player := range playersInRoom {
		playerInterfaces[i] = &player
	}
	roomInArea.SetPlayers(playerInterfaces)
}

func setItems(db *sql.DB, roomInArea interfaces.Room) {
	// retrieve the items and attach them to the room in the WorldState
	itemsInRoom, err := items.GetItemsInRoom(db, roomInArea.GetUUID())
	if err != nil {
		fmt.Printf("error retrieving items: %v", err)
	}

	itemInterfaces := make([]interfaces.Item, len(itemsInRoom))
	for i, item := range itemsInRoom {
		itemInterfaces[i] = &item
	}
	roomInArea.SetItems(itemInterfaces)
}

package world_state

import (
	"database/sql"
	"fmt"
	"mud/areas"
	"mud/interfaces"
	"mud/items"
	"mud/players"
)

type WorldState struct {
	Areas         map[string]*areas.Area
	RoomToAreaMap map[string]string
	DB            *sql.DB
}

func NewWorldState(areas map[string]*areas.Area, roomToAreaMap map[string]string, db *sql.DB) *WorldState {
	return &WorldState{Areas: areas, RoomToAreaMap: roomToAreaMap, DB: db}
}

func getRoomFromAreaByUUID(area *areas.Area, roomUUID string) *areas.Room {
	for _, room := range area.Rooms {
		if room.UUID == roomUUID {
			return &room
		}
	}
	fmt.Println("Room not found")
	return nil
}

func (worldState *WorldState) RemovePlayerFromRoom(roomUUID string, player interfaces.Player) {
	areaUUID := worldState.RoomToAreaMap[roomUUID]
	roomIdx := -1

	for idx, roomInArea := range worldState.Areas[areaUUID].Rooms {
		if roomInArea.UUID == roomUUID {
			roomIdx = idx
			break
		}
	}

	if roomIdx >= 0 {
		playerIdx := -1
		for idx, playerInRoom := range worldState.Areas[areaUUID].Rooms[roomIdx].Players {
			if playerInRoom.GetUUID() == player.GetUUID() {
				playerIdx = idx
			}
		}
		if playerIdx >= 0 {
			playersInRoom := worldState.Areas[areaUUID].Rooms[roomIdx].Players
			worldState.Areas[areaUUID].Rooms[roomIdx].Players = append(playersInRoom[:playerIdx], playersInRoom[playerIdx+1:]...)
		}

	}
}

func (worldState *WorldState) AddPlayerToRoom(roomUUID string, player interfaces.Player) {
	areaUUID := worldState.RoomToAreaMap[roomUUID]
	roomIdx := -1
	for idx, roomInArea := range worldState.Areas[areaUUID].Rooms {
		if roomInArea.UUID == roomUUID {
			roomIdx = idx
			break
		}
	}

	if roomIdx >= 0 {
		playersInRoom := worldState.Areas[areaUUID].Rooms[roomIdx].Players
		worldState.Areas[areaUUID].Rooms[roomIdx].Players = append(playersInRoom, player)
	}
}

func (worldState *WorldState) TransferItemFromRoomToPlayer(room *areas.Room, item interfaces.Item, player interfaces.Player) error {
	areaUUID, ok := worldState.RoomToAreaMap[room.UUID]
	if !ok {
		return fmt.Errorf("area UUID not found for room UUID: %s", room.UUID)
	}
	roomIndex := -1
	for idx, roomInArea := range worldState.Areas[areaUUID].Rooms {
		if roomInArea.UUID == room.UUID {
			roomIndex = idx
			break
		}
	}
	itemIndex := -1
	for idx := range room.Items {
		if room.Items[idx] == item {
			itemIndex = idx
			break
		}
	}
	if itemIndex > -1 && roomIndex > -1 {
		updatedItems := append(room.Items[:itemIndex], room.Items[itemIndex+1:]...)
		worldState.Areas[areaUUID].Rooms[roomIndex].Items = updatedItems
	}

	player.AddItemToInventory(worldState.DB, item)

	return nil
}

func (worldState *WorldState) GetRoom(roomUUID string, followExits bool) *areas.Room {
	var retrievedRoom *areas.Room
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
		for _, roomInArea := range area.Rooms {
			if roomInArea.UUID == roomUUID {
				retrievedRoom = &roomInArea
				break
			}
		}
	} else {
		fmt.Println("rooms in memory does not match rooms in DB, must make a call to the DB")
		queryString := `
			SELECT r.UUID, r.area_uuid, r.name, r.description,
				r.exit_north, r.exit_south, r.exit_east, r.exit_west,
				r.exit_up, r.exit_down,
				a.UUID AS area_uuid, a.name AS area_name, a.description AS area_description
			FROM rooms r
			LEFT JOIN areas a ON r.area_uuid = a.UUID
			WHERE a.UUID = ?`
		rows, err := worldState.DB.Query(queryString, areaUUID)
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
			area.Rooms = append(area.Rooms, *room)
			if room.UUID == roomUUID {
				retrievedRoom = room
			}
		}
		defer rows.Close()

		// now that we have loaded all the rooms in the area, we can go back and hook up all the
		// exits.  if one of the exits happens to exist in a different area, we can make a db query to retrieve that one.
		for idx, roomInArea := range area.Rooms {
			if followExits {
				if roomInArea.Exits.South != nil {
					if roomInArea.Exits.South.Name == "" {
						exitRoom := getRoomFromAreaByUUID(area, roomInArea.Exits.South.UUID)
						if exitRoom == nil {
							rm, err := areas.GetRoomFromDB(roomInArea.Exits.South.UUID, worldState.DB)
							if err != nil {
								fmt.Printf("error getting room: %v", err)
							}
							exitRoom = rm
						}
						roomInArea.Exits.South = exitRoom
					}
				}
				if roomInArea.Exits.North != nil {
					if roomInArea.Exits.North.Name == "" {
						exitRoom := getRoomFromAreaByUUID(area, roomInArea.Exits.North.UUID)
						if exitRoom == nil {
							rm, err := areas.GetRoomFromDB(roomInArea.Exits.North.UUID, worldState.DB)
							if err != nil {
								fmt.Printf("error getting room: %v", err)
							}
							exitRoom = rm
						}
						roomInArea.Exits.North = exitRoom
					}
				}
				if roomInArea.Exits.West != nil {
					if roomInArea.Exits.West.Name == "" {
						exitRoom := getRoomFromAreaByUUID(area, roomInArea.Exits.West.UUID)
						if exitRoom == nil {
							rm, err := areas.GetRoomFromDB(roomInArea.Exits.West.UUID, worldState.DB)
							if err != nil {
								fmt.Printf("error getting room: %v", err)
							}
							exitRoom = rm
						}
						roomInArea.Exits.West = exitRoom
					}
				}
				if roomInArea.Exits.East != nil {
					if roomInArea.Exits.East.Name == "" {
						exitRoom := getRoomFromAreaByUUID(area, roomInArea.Exits.East.UUID)
						if exitRoom == nil {
							rm, err := areas.GetRoomFromDB(roomInArea.Exits.East.UUID, worldState.DB)
							if err != nil {
								fmt.Printf("error getting room: %v", err)
							}
							exitRoom = rm
						}
						roomInArea.Exits.East = exitRoom
					}
				}
				if roomInArea.Exits.Up != nil {
					if roomInArea.Exits.Up.Name == "" {
						exitRoom := getRoomFromAreaByUUID(area, roomInArea.Exits.Up.UUID)
						if exitRoom == nil {
							rm, err := areas.GetRoomFromDB(roomInArea.Exits.Up.UUID, worldState.DB)
							if err != nil {
								fmt.Printf("error getting room: %v", err)
							}
							exitRoom = rm
						}
						roomInArea.Exits.Up = exitRoom
					}
				}
				if roomInArea.Exits.Down != nil {
					if roomInArea.Exits.Down.Name == "" {
						exitRoom := getRoomFromAreaByUUID(area, roomInArea.Exits.Down.UUID)
						if exitRoom == nil {
							rm, err := areas.GetRoomFromDB(roomInArea.Exits.Down.UUID, worldState.DB)
							if err != nil {
								fmt.Printf("error getting room: %v", err)
							}
							exitRoom = rm
						}
						roomInArea.Exits.Down = exitRoom
					}
				}
			}
			itemsInRoom, err := items.GetItemsInRoom(worldState.DB, roomInArea.UUID)
			if err != nil {
				fmt.Printf("error retrieving items: %v", err)
			}
			itemInterfaces := make([]interfaces.Item, len(itemsInRoom))
			for i, item := range itemsInRoom {
				itemInterfaces[i] = &item
			}
			roomInArea.Items = itemInterfaces

			// retrieve the players and attach them to the room
			playersInRoom, err := players.GetPlayersInRoom(worldState.DB, roomInArea.UUID)
			if err != nil {
				fmt.Printf("error retrieving players: %v", err)
			}
			playerInterfaces := make([]interfaces.Player, len(playersInRoom))
			for i, player := range playersInRoom {
				playerInterfaces[i] = &player
			}

			roomPlayerInterfaces := make([]areas.PlayerInRoomInterface, len(playerInterfaces))
			for i, playerInterface := range playerInterfaces {
				roomPlayerInterfaces[i] = playerInterface
			}
			roomInArea.Players = roomPlayerInterfaces

			area.Rooms[idx] = roomInArea
		}
	}

	return retrievedRoom
}

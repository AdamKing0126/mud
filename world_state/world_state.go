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

func findIndex(slice interface{}, matchFunc func(int) bool) int {
	switch s := slice.(type) {
	case []areas.Room:
		for idx := range s {
			if matchFunc(idx) {
				return idx
			}
		}
	case []interfaces.Player:
		for idx := range s {
			if matchFunc(idx) {
				return idx
			}
		}
	case []areas.PlayerInRoomInterface:
		for idx := range s {
			if matchFunc(idx) {
				return idx
			}
		}
	case []interfaces.Item:
		for idx := range s {
			if matchFunc(idx) {
				return idx
			}
		}
	}
	return -1
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

func (worldState *WorldState) RemovePlayerFromRoom(roomUUID string, player interfaces.Player) error {
	areaUUID := worldState.RoomToAreaMap[roomUUID]
	area := worldState.Areas[areaUUID]
	rooms := area.Rooms

	roomIdx := findIndex(area.Rooms, func(idx int) bool { return area.Rooms[idx].UUID == roomUUID })
	if roomIdx == -1 {
		return fmt.Errorf("room with UUID %s not found in area %s", roomUUID, areaUUID)
	}
	room := rooms[roomIdx]

	playerIdx := findIndex(room.Players, func(idx int) bool { return room.Players[idx] == player })
	if playerIdx == -1 {
		return fmt.Errorf("player with UUID %s not found in room %s", player.GetUUID(), roomUUID)
	}

	playersInRoom := worldState.Areas[areaUUID].Rooms[roomIdx].Players
	room.Players = append(playersInRoom[:playerIdx], playersInRoom[playerIdx+1:]...)

	worldState.Areas[areaUUID].Rooms = append(area.Rooms[:roomIdx], append([]areas.Room{room}, area.Rooms[roomIdx+1:]...)...)
	return nil
}

func (worldState *WorldState) AddPlayerToRoom(roomUUID string, player interfaces.Player) error {
	areaUUID := worldState.RoomToAreaMap[roomUUID]
	area := worldState.Areas[areaUUID]
	roomIdx := findIndex(area.Rooms, func(idx int) bool { return area.Rooms[idx].UUID == roomUUID })
	if roomIdx == -1 {
		return fmt.Errorf("room UUID %s not found in area %s", roomUUID, areaUUID)
	}

	playersInRoom := worldState.Areas[areaUUID].Rooms[roomIdx].Players
	worldState.Areas[areaUUID].Rooms[roomIdx].Players = append(playersInRoom, player)
	return nil
}

func (worldState *WorldState) TransferItemFromRoomToPlayer(room *areas.Room, item interfaces.Item, player interfaces.Player) error {
	areaUUID, ok := worldState.RoomToAreaMap[room.UUID]
	if !ok {
		return fmt.Errorf("area UUID not found for room UUID: %s", room.UUID)
	}
	area := worldState.Areas[areaUUID]

	roomIndex := findIndex(area.Rooms, func(idx int) bool { return area.Rooms[idx].UUID == room.UUID })
	if roomIndex == -1 {
		return fmt.Errorf("room UUID %s not found in area %s", room.UUID, areaUUID)
	}

	playerIndex := findIndex(room.Players, func(idx int) bool { return room.Players[idx].GetUUID() == player.GetUUID() })
	if playerIndex == -1 {
		return fmt.Errorf("player UUID %s not found in room %s", player.GetUUID(), room.UUID)
	}

	itemIndex := findIndex(room.Items, func(idx int) bool { return room.Items[idx] == item })
	if itemIndex == -1 {
		return fmt.Errorf("item %s not found in room %s", item.GetUUID(), room.UUID)
	}

	updatedItems := append(room.Items[:itemIndex], room.Items[itemIndex+1:]...)
	room.Items = updatedItems
	player.AddItemToInventory(worldState.DB, item)
	room.Players = append(room.Players[:playerIndex], append([]areas.PlayerInRoomInterface{player}, room.Players[playerIndex+1:]...)...)
	worldState.Areas[areaUUID].Rooms = append(area.Rooms[:roomIndex], append([]areas.Room{*room}, area.Rooms[roomIndex+1:]...)...)

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

func (worldState *WorldState) RemoveItemFromPlayerInventory(player interfaces.Player, item interfaces.Item) error {
	itemIdx := findIndex(player.GetInventory(), func(idx int) bool { return player.GetInventory()[idx] == item })
	if itemIdx == -1 {
		return fmt.Errorf("no item %s found in inventory for player %s", item.GetUUID(), player.GetUUID())
	}

	area := worldState.Areas[player.GetArea()]
	roomIdx := findIndex(area.Rooms, func(idx int) bool { return area.Rooms[idx].UUID == player.GetRoomUUID() })
	if roomIdx == -1 {
		return fmt.Errorf("no room %s found in area %s", player.GetRoomUUID(), area.UUID)
	}

	room := area.Rooms[roomIdx]
	playerIndex := findIndex(room.Players, func(idx int) bool { return room.Players[idx].GetUUID() == player.GetUUID() })
	if playerIndex == -1 {
		return fmt.Errorf("no player %s found in room %s", player.GetUUID(), room.UUID)
	}

	playerInventory := player.GetInventory()
	itemIndex := findIndex(playerInventory, func(idx int) bool { return playerInventory[idx].GetUUID() == item.GetUUID() })
	if itemIndex == -1 {
		return fmt.Errorf("no item %s found in player %s inventory", item.GetUUID(), player.GetUUID())
	}

	// TODO ADAM FIX
	// set the inventory for the slice excluding the itemIndex
	// set the room's players, replacing room.Players[playerIndex] with player
	fmt.Printf("inventory: %d", len(playerInventory))

	return nil
}

func (worldState *WorldState) TransferItemFromPlayerToPlayer(item interfaces.Item, player interfaces.Player, receiverName string) (areas.PlayerInRoomInterface, areas.PlayerInRoomInterface, error) {
	areaUUID := player.GetArea()
	roomUUID := player.GetRoomUUID()
	area := worldState.Areas[areaUUID]
	roomIndex := findIndex(area.Rooms, func(idx int) bool { return area.Rooms[idx].UUID == roomUUID })
	if roomIndex < 0 {
		return nil, nil, fmt.Errorf("room %s not found in area %s", roomUUID, areaUUID)
	}

	room := area.Rooms[roomIndex]
	giverIndex := findIndex(room.Players, func(idx int) bool { return room.Players[idx].GetUUID() == player.GetUUID() })
	if giverIndex < 0 {
		return nil, nil, fmt.Errorf("player %s not found in room %s", player.GetUUID(), roomUUID)
	}
	giver := room.Players[giverIndex]

	receiverIndex := findIndex(room.Players, func(idx int) bool { return room.Players[idx].GetName() == receiverName })
	if receiverIndex < 0 {
		return nil, nil, fmt.Errorf("player %s not found in room %s", receiverName, roomUUID)
	}
	receiver := room.Players[receiverIndex]

	itemIndex := findIndex(giver.GetInventory(), func(idx int) bool { return giver.GetInventory()[idx].GetUUID() == item.GetUUID() })
	if itemIndex < 0 {
		return nil, nil, fmt.Errorf("item %s not found for player inventory %s", item.GetUUID(), giver.GetUUID())
	}
	giver.SetInventory(append(giver.GetInventory()[:itemIndex], giver.GetInventory()[itemIndex+1:]...))
	receiver.SetInventory(append(receiver.GetInventory(), item))
	return giver, receiver, nil
}

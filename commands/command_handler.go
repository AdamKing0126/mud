package commands

import (
	"database/sql"
	"fmt"
	"mud/areas"
	"mud/display"
	"mud/interfaces"
	"mud/items"
	"mud/utils"
)

var CommandHandlers = map[string]utils.CommandHandlerWithPriority{
	"north":     {Handler: &MovePlayerCommandHandler{Direction: "north"}, Priority: 1},
	"south":     {Handler: &MovePlayerCommandHandler{Direction: "south"}, Priority: 1},
	"west":      {Handler: &MovePlayerCommandHandler{Direction: "west"}, Priority: 1},
	"east":      {Handler: &MovePlayerCommandHandler{Direction: "east"}, Priority: 1},
	"up":        {Handler: &MovePlayerCommandHandler{Direction: "up"}, Priority: 1},
	"down":      {Handler: &MovePlayerCommandHandler{Direction: "down"}, Priority: 1},
	"look":      {Handler: &LookCommandHandler{}, Priority: 2},
	"logout":    {Handler: &LogoutCommandHandler{}, Priority: 10},
	"exits":     {Handler: &ExitsCommandHandler{}, Priority: 2},
	"take":      {Handler: &TakeCommandHandler{}, Priority: 2},
	"drop":      {Handler: &DropCommandHandler{}, Priority: 2},
	"inventory": {Handler: &InventoryCommandHandler{}, Priority: 2},
	"foo":       {Handler: &FooCommandHandler{}, Priority: 2},
}

func getRoom(roomUUID string, db *sql.DB) (*areas.Room, error) {
	query := `
		SELECT r.UUID, r.area_uuid, r.name, r.description,
			r.exit_north, r.exit_south, r.exit_east, r.exit_west,
			r.exit_up, r.exit_down,
			a.UUID AS area_uuid, a.name AS area_name, a.description AS area_description
		FROM rooms r
		LEFT JOIN areas a ON r.area_uuid = a.UUID
		WHERE r.UUID = ?`

	room_rows, err := db.Query(query, roomUUID)
	if err != nil {
		return nil, err
	}

	defer room_rows.Close()
	if !room_rows.Next() {
		return nil, fmt.Errorf("room with UUID %s does not exist", roomUUID)
	}

	var northExitUUID, southExitUUID, eastExitUUID, westExitUUID, upExitUUID, downExitUUID string
	room := &areas.Room{Exits: areas.ExitInfo{}}
	err = room_rows.Scan(
		&room.UUID, &room.AreaUUID, &room.Name, &room.Description,
		&northExitUUID, &southExitUUID, &eastExitUUID, &westExitUUID,
		&upExitUUID, &downExitUUID,
		&room.Area.UUID, &room.Area.Name, &room.Area.Description,
	)
	if err != nil {
		return nil, err
	}

	if northExitUUID != "" {
		room.Exits.North = &areas.Room{UUID: northExitUUID}
	}

	if southExitUUID != "" {
		room.Exits.South = &areas.Room{UUID: southExitUUID}
	}

	if westExitUUID != "" {
		room.Exits.West = &areas.Room{UUID: westExitUUID}
	}

	if eastExitUUID != "" {
		room.Exits.East = &areas.Room{UUID: eastExitUUID}
	}

	if downExitUUID != "" {
		room.Exits.Down = &areas.Room{UUID: downExitUUID}
	}

	if upExitUUID != "" {
		room.Exits.Up = &areas.Room{UUID: upExitUUID}
	}

	items, err := items.GetItemsInRoom(db, roomUUID)
	if err != nil {
		return nil, err
	}

	itemInterfaces := make([]interfaces.ItemInterface, len(items))
	copy(itemInterfaces, items)

	room.Items = itemInterfaces

	return room, nil
}

type MovePlayerCommandHandler struct {
	Direction   string
	LookHandler *LookCommandHandler
}

func (h *MovePlayerCommandHandler) Execute(db *sql.DB, player interfaces.PlayerInterface, command string, arguments []string, currentChannel chan interfaces.ActionInterface, updateChannel func(string)) {
	roomUUID := player.GetRoom()
	areaUUID := player.GetArea()
	playerConn := player.GetConn()
	currentRoom, err := getRoom(roomUUID, db)
	if err != nil {
		fmt.Fprintf(playerConn, "%v", err)
	}
	switch h.Direction {
	case "north":
		if currentRoom.Exits.North == nil {
			fmt.Fprintf(playerConn, "You cannot go that way.\n")

		} else {
			fmt.Fprintf(playerConn, "=======================\n\n")
			player.SetLocation(db, currentRoom.Exits.North.UUID)
			var lookArgs []string
			h.LookHandler.Execute(db, player, "look", lookArgs, currentChannel, updateChannel)
		}
	case "south":
		if currentRoom.Exits.South == nil {
			fmt.Fprintf(playerConn, "You cannot go that way.\n")

		} else {
			fmt.Fprintf(playerConn, "=======================\n\n")
			player.SetLocation(db, currentRoom.Exits.South.UUID)
			var lookArgs []string
			h.LookHandler.Execute(db, player, "look", lookArgs, currentChannel, updateChannel)
		}
	case "west":
		if currentRoom.Exits.West == nil {
			fmt.Fprintf(playerConn, "You cannot go that way.\n")

		} else {
			fmt.Fprintf(playerConn, "=======================\n\n")
			player.SetLocation(db, currentRoom.Exits.West.UUID)
			var lookArgs []string
			h.LookHandler.Execute(db, player, "look", lookArgs, currentChannel, updateChannel)
		}
	case "east":
		if currentRoom.Exits.East == nil {
			fmt.Fprintf(playerConn, "You cannot go that way.\n")

		} else {
			fmt.Fprintf(playerConn, "=======================\n\n")
			player.SetLocation(db, currentRoom.Exits.East.UUID)
			var lookArgs []string
			h.LookHandler.Execute(db, player, "look", lookArgs, currentChannel, updateChannel)
		}
	case "up":
		if currentRoom.Exits.Up == nil {
			fmt.Fprintf(playerConn, "You cannot go that way.\n")

		} else {
			fmt.Fprintf(playerConn, "=======================\n\n")
			player.SetLocation(db, currentRoom.Exits.Up.UUID)
			var lookArgs []string
			h.LookHandler.Execute(db, player, "look", lookArgs, currentChannel, updateChannel)
		}
	default:
		if currentRoom.Exits.Down == nil {
			fmt.Fprintf(playerConn, "You cannot go that way.\n")
		} else {
			fmt.Fprintf(playerConn, "=======================\n\n")
			player.SetLocation(db, currentRoom.Exits.Down.UUID)
			var lookArgs []string
			h.LookHandler.Execute(db, player, "look", lookArgs, currentChannel, updateChannel)
		}
	}

	if areaUUID != player.GetArea() {
		updateChannel(player.GetArea())
	}
}

type ExitsCommandHandler struct{}

func (h *ExitsCommandHandler) Execute(db *sql.DB, player interfaces.PlayerInterface, command string, arguments []string, currentChannel chan interfaces.ActionInterface, updateChannel func(string)) {
	roomUUID := player.GetRoom()
	playerConn := player.GetConn()
	currentRoom, err := getRoom(roomUUID, db)
	if err != nil {
		fmt.Fprintf(playerConn, "%v", err)
	}
	if currentRoom.Exits.North != nil {
		northExit, err := getRoom(currentRoom.Exits.North.UUID, db)
		if err != nil {
			fmt.Fprintf(playerConn, "%v", err)
		}
		fmt.Fprintf(playerConn, "North: %s\n", northExit.Name)
	}
	if currentRoom.Exits.South != nil {
		southExit, err := getRoom(currentRoom.Exits.South.UUID, db)
		if err != nil {
			fmt.Fprintf(playerConn, "%v", err)
		}
		fmt.Fprintf(playerConn, "South: %s\n", southExit.Name)
	}
	if currentRoom.Exits.West != nil {
		westExit, err := getRoom(currentRoom.Exits.West.UUID, db)
		if err != nil {
			fmt.Fprintf(playerConn, "%v", err)
		}
		fmt.Fprintf(playerConn, "West: %s\n", westExit.Name)
	}
	if currentRoom.Exits.East != nil {
		eastExit, err := getRoom(currentRoom.Exits.East.UUID, db)
		if err != nil {
			fmt.Fprintf(playerConn, "%v", err)
		}
		fmt.Fprintf(playerConn, "East: %s\n", eastExit.Name)
	}
	if currentRoom.Exits.Up != nil {
		upExit, err := getRoom(currentRoom.Exits.Up.UUID, db)
		if err != nil {
			fmt.Fprintf(playerConn, "%v", err)
		}
		fmt.Fprintf(playerConn, "Up: %s\n", upExit.Name)
	}
	if currentRoom.Exits.Down != nil {
		downExit, err := getRoom(currentRoom.Exits.Down.UUID, db)
		if err != nil {
			fmt.Fprintf(playerConn, "%v", err)
		}
		fmt.Fprintf(playerConn, "Down: %s\n", downExit.Name)
	}
}

type LogoutCommandHandler struct{}

func (h *LogoutCommandHandler) Execute(db *sql.DB, player interfaces.PlayerInterface, command string, arguments []string, currentChannel chan interfaces.ActionInterface, updateChannel func(string)) {
	playerConn := player.GetConn()
	fmt.Fprintf(playerConn, "Goodbye!\n")
	player.Logout()
}

type LookCommandHandler struct{}

func (h *LookCommandHandler) Execute(db *sql.DB, player interfaces.PlayerInterface, command string, arguments []string, currentChannel chan interfaces.ActionInterface, updateChannel func(string)) {
	roomUUID := player.GetRoom()
	playerConn := player.GetConn()
	currentRoom, err := getRoom(roomUUID, db)
	if err != nil {
		fmt.Fprintf(playerConn, "%v", err)
	}

	if len(arguments) == 0 {

		display.PrintWithColor(player, fmt.Sprintf("%s\n", currentRoom.Area.Name), "primary")
		display.PrintWithColor(player, fmt.Sprintf("%s\n", currentRoom.Area.Description), "secondary")
		display.PrintWithColor(player, "-----------------------\n\n", "secondary")
		fmt.Fprintf(playerConn, "%s\n", currentRoom.Name)
		fmt.Fprintf(playerConn, "%s\n", currentRoom.Description)

		if len(currentRoom.Items) > 0 {
			fmt.Fprintf(playerConn, "You see the following items:\n")
			for _, item := range currentRoom.Items {
				fmt.Fprintf(playerConn, "%s\n", item.GetName())
			}
		}
	} else if len(arguments) == 1 {
		switch arguments[0] {
		case "north":
			if currentRoom.Exits.North != nil {
				northExit, err := getRoom(currentRoom.Exits.North.UUID, db)
				if err != nil {
					fmt.Fprintf(playerConn, "%v", err)
				}
				fmt.Fprintf(playerConn, "You look North. You see %s\n", northExit.Name)
			} else {
				fmt.Fprintf(playerConn, "You don't see anything in that direction\n")
			}
		case "south":
			if currentRoom.Exits.South != nil {
				southExit, err := getRoom(currentRoom.Exits.South.UUID, db)
				if err != nil {
					fmt.Fprintf(playerConn, "%v", err)
				}
				fmt.Fprintf(playerConn, "You look South. You see %s\n", southExit.Name)
			} else {
				fmt.Fprintf(playerConn, "You don't see anything in that direction\n")
			}
		case "east":
			if currentRoom.Exits.East != nil {
				eastExit, err := getRoom(currentRoom.Exits.East.UUID, db)
				if err != nil {
					fmt.Fprintf(playerConn, "%v", err)
				}
				fmt.Fprintf(playerConn, "You look East. You see %s\n", eastExit.Name)
			} else {
				fmt.Fprintf(playerConn, "You don't see anything in that direction\n")
			}
		case "west":
			if currentRoom.Exits.West != nil {
				westExit, err := getRoom(currentRoom.Exits.West.UUID, db)
				if err != nil {
					fmt.Fprintf(playerConn, "%v", err)
				}
				fmt.Fprintf(playerConn, "You look West. You see %s\n", westExit.Name)
			} else {
				fmt.Fprintf(playerConn, "You don't see anything in that direction\n")
			}
		case "up":
			if currentRoom.Exits.Up != nil {
				upExit, err := getRoom(currentRoom.Exits.Up.UUID, db)
				if err != nil {
					fmt.Fprintf(playerConn, "%v", err)
				}
				fmt.Fprintf(playerConn, "You look Up. You see %s\n", upExit.Name)
			} else {
				fmt.Fprintf(playerConn, "You don't see anything in that direction\n")
			}
		case "down":
			if currentRoom.Exits.Down != nil {
				downExit, err := getRoom(currentRoom.Exits.Down.UUID, db)
				if err != nil {
					fmt.Fprintf(playerConn, "%v", err)
				}
				fmt.Fprintf(playerConn, "You look Down. You see %s\n", downExit.Name)
			} else {
				fmt.Fprintf(playerConn, "You don't see anything in that direction\n")
			}
		default:
			itemName := arguments[0]
			found := false
			itemsForPlayer, err := items.GetItemsForPlayer(db, player.GetUUID())
			if err != nil {
				fmt.Fprintf(playerConn, "%v", err)
			}

			items := append(currentRoom.Items, itemsForPlayer...)
			for _, item := range items {
				if item.GetName() == itemName {
					fmt.Fprintf(playerConn, "%s\n", item.GetDescription())
					found = true
					break
				}
			}

			if !found {
				fmt.Fprintf(playerConn, "You don't see that.\n")
			}
		}
	} else {
		fmt.Fprintf(playerConn, "I don't know how to do that yet.\n")
	}
}

type TakeCommandHandler struct{}

func (h *TakeCommandHandler) Execute(db *sql.DB, player interfaces.PlayerInterface, command string, arguments []string, currentChannel chan interfaces.ActionInterface, updateChannel func(string)) {
	roomUUID := player.GetRoom()
	playerConn := player.GetConn()
	currentRoom, err := getRoom(roomUUID, db)
	if err != nil {
		fmt.Fprintf(playerConn, "%v", err)
	}

	if len(currentRoom.Items) > 0 {
		for _, item := range currentRoom.Items {
			if item.GetName() == arguments[0] {
				query := "UPDATE item_locations SET room_uuid = '', player_uuid = ? WHERE item_uuid = ?"
				_, err := db.Exec(query, player.GetUUID(), item.GetUUID())
				if err != nil {
					fmt.Fprintf(playerConn, "Failed to update item location: %v\n", err)
				}
				fmt.Fprintf(playerConn, "You take the %s.\n", item.GetName())
				break
			}
		}
	} else {
		fmt.Fprintf(playerConn, "You don't see that here.\n")
	}
}

type DropCommandHandler struct{}

func (h *DropCommandHandler) Execute(db *sql.DB, player interfaces.PlayerInterface, command string, arguments []string, currentChannel chan interfaces.ActionInterface, updateChannel func(string)) {
	roomUUID := player.GetRoom()
	playerConn := player.GetConn()

	playerItems, err := items.GetItemsForPlayer(db, player.GetUUID())
	if err != nil {
		fmt.Fprintf(playerConn, "%v", err)
	}

	if len(playerItems) > 0 {
		for _, item := range playerItems {
			if item.GetName() == arguments[0] {
				query := "UPDATE item_locations SET room_uuid = ?, player_uuid = NULL WHERE item_uuid = ?"
				_, err := db.Exec(query, roomUUID, item.GetUUID())
				if err != nil {
					fmt.Fprintf(playerConn, "Failed to update item location: %v\n", err)
				}
				fmt.Fprintf(playerConn, "You drop the %s.\n", item.GetName())
				break
			}
		}
	} else {
		fmt.Fprintf(playerConn, "You don't have that item.\n")
	}
}

type InventoryCommandHandler struct{}

func (h *InventoryCommandHandler) Execute(db *sql.DB, player interfaces.PlayerInterface, command string, arguments []string, currentChannel chan interfaces.ActionInterface, updateChannel func(string)) {
	playerConn := player.GetConn()
	playerItems, err := items.GetItemsForPlayer(db, player.GetUUID())
	if err != nil {
		fmt.Fprintf(playerConn, "%v", err)
	}

	fmt.Fprintf(playerConn, "You are carrying:\n")

	if len(playerItems) == 0 {
		fmt.Fprintf(playerConn, "Nothing\n")
	} else {
		for _, item := range playerItems {
			fmt.Fprintf(playerConn, "%s\n", item.GetName())
		}
	}
}

type FooCommandHandler struct{}

func (h *FooCommandHandler) Execute(db *sql.DB, player interfaces.PlayerInterface, command string, arguments []string, currentChannel chan interfaces.ActionInterface, updateChannel func(string)) {
	currentChannel <- &areas.Action{Player: player, Command: command, Arguments: arguments}
}

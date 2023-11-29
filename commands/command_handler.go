package commands

import (
	"database/sql"
	"fmt"
	"mud/areas"
	"mud/interfaces"
	"mud/items"
	"strings"
)

type CommandParser struct {
	commandName string
	arguments   []string
}

type CommandHandler func(*sql.DB, interfaces.PlayerInterface, string, []string, chan interfaces.ActionInterface, func(string)) interface{}

func NewCommandParser(commandString string) *CommandParser {
	commandString = strings.ToLower(commandString)

	// Split the command string into the command name and arguments.
	commandParts := strings.Split(commandString, " ")
	commandName := commandParts[0]
	arguments := commandParts[1:]
	// look for partial commands ie "n" for "north"
	for fullCommand := range CommandHandlers {
		if strings.HasPrefix(strings.ToLower(fullCommand), commandParts[0]) {
			return &CommandParser{
				commandName: fullCommand,
				arguments:   arguments,
			}
		}
	}

	return &CommandParser{
		commandName: commandName,
		arguments:   arguments,
	}
}

func (p *CommandParser) GetCommandName() string {
	return p.commandName
}

func (p *CommandParser) GetArguments() []string {
	return p.arguments
}

var CommandHandlers map[string]CommandHandler = map[string]CommandHandler{
	"north":     curryMovePlayerCommand("north"),
	"south":     curryMovePlayerCommand("south"),
	"west":      curryMovePlayerCommand("west"),
	"east":      curryMovePlayerCommand("east"),
	"up":        curryMovePlayerCommand("up"),
	"down":      curryMovePlayerCommand("down"),
	"look":      HandleLookCommand,
	"logout":    HandleLogoutCommand,
	"exits":     HandleExitsCommand,
	"take":      HandleTakeCommand,
	"drop":      HandleDropCommand,
	"inventory": HandleInventoryCommand,
	"foo":       QueueCommandHandler(HandleFooCommand),
}

func QueueCommandHandler(handler CommandHandler) CommandHandler {
	return func(db *sql.DB, player interfaces.PlayerInterface, command string, arguments []string, currentChannel chan interfaces.ActionInterface, updateChannel func(string)) interface{} {
		return func() {
			fmt.Println("Queueing command: ", command)
			currentChannel <- &areas.Action{Player: player, Command: command}
			if function, ok := handler(db, player, command, arguments, currentChannel, updateChannel).(func()); ok {
				function()
			}
		}
	}
}

func curryMovePlayerCommand(direction string) CommandHandler {
	return func(db *sql.DB, player interfaces.PlayerInterface, commandName string, arguments []string, currentChannel chan interfaces.ActionInterface, updateChannel func(string)) interface{} {
		HandleMovePlayerCommand(db, player, commandName, arguments, currentChannel, updateChannel)
		return nil
	}
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

func HandleMovePlayerCommand(db *sql.DB, player interfaces.PlayerInterface, command string, arguments []string, currentChannel chan interfaces.ActionInterface, updateChannel func(string)) {
	roomUUID := player.GetRoom()
	areaUUID := player.GetArea()
	playerConn := player.GetConn()
	currentRoom, err := getRoom(roomUUID, db)
	if err != nil {
		fmt.Fprintf(playerConn, "%v", err)
		return
	}
	switch command {
	case "north":
		if currentRoom.Exits.North == nil {
			fmt.Fprintf(playerConn, "You cannot go that way.\n")

		} else {
			fmt.Fprintf(playerConn, "=======================\n\n")
			player.SetLocation(db, currentRoom.Exits.North.UUID)
			var lookArgs []string
			HandleLookCommand(db, player, "look", lookArgs, currentChannel, updateChannel)
		}
	case "south":
		if currentRoom.Exits.South == nil {
			fmt.Fprintf(playerConn, "You cannot go that way.\n")

		} else {
			fmt.Fprintf(playerConn, "=======================\n\n")
			player.SetLocation(db, currentRoom.Exits.South.UUID)
			var lookArgs []string
			HandleLookCommand(db, player, "look", lookArgs, currentChannel, updateChannel)
		}
	case "west":
		if currentRoom.Exits.West == nil {
			fmt.Fprintf(playerConn, "You cannot go that way.\n")

		} else {
			fmt.Fprintf(playerConn, "=======================\n\n")
			player.SetLocation(db, currentRoom.Exits.West.UUID)
			var lookArgs []string
			HandleLookCommand(db, player, "look", lookArgs, currentChannel, updateChannel)
		}
	case "east":
		if currentRoom.Exits.East == nil {
			fmt.Fprintf(playerConn, "You cannot go that way.\n")

		} else {
			fmt.Fprintf(playerConn, "=======================\n\n")
			player.SetLocation(db, currentRoom.Exits.East.UUID)
			var lookArgs []string
			HandleLookCommand(db, player, "look", lookArgs, currentChannel, updateChannel)
		}
	case "up":
		if currentRoom.Exits.Up == nil {
			fmt.Fprintf(playerConn, "You cannot go that way.\n")

		} else {
			fmt.Fprintf(playerConn, "=======================\n\n")
			player.SetLocation(db, currentRoom.Exits.Up.UUID)
			var lookArgs []string
			HandleLookCommand(db, player, "look", lookArgs, currentChannel, updateChannel)
		}
	default:
		if currentRoom.Exits.Down == nil {
			fmt.Fprintf(playerConn, "You cannot go that way.\n")
		} else {
			fmt.Fprintf(playerConn, "=======================\n\n")
			player.SetLocation(db, currentRoom.Exits.Down.UUID)
			var lookArgs []string
			HandleLookCommand(db, player, "look", lookArgs, currentChannel, updateChannel)
		}
	}

	if areaUUID != player.GetArea() {
		updateChannel(player.GetArea())
	}
}

func HandleExitsCommand(db *sql.DB, player interfaces.PlayerInterface, command string, arguments []string, currentChannel chan interfaces.ActionInterface, updateChannel func(string)) interface{} {
	roomUUID := player.GetRoom()
	playerConn := player.GetConn()
	currentRoom, err := getRoom(roomUUID, db)
	if err != nil {
		fmt.Fprintf(playerConn, "%v", err)
		return nil
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

	return nil
}

func HandleLogoutCommand(db *sql.DB, player interfaces.PlayerInterface, command string, arguments []string, currentChannel chan interfaces.ActionInterface, updateChannel func(string)) interface{} {
	playerConn := player.GetConn()
	fmt.Fprintf(playerConn, "Goodbye!\n")
	player.Logout()
	return nil
}

func HandleLookCommand(db *sql.DB, player interfaces.PlayerInterface, command string, arguments []string, currentChannel chan interfaces.ActionInterface, updateChannel func(string)) interface{} {
	roomUUID := player.GetRoom()
	playerConn := player.GetConn()
	currentRoom, err := getRoom(roomUUID, db)
	if err != nil {
		fmt.Fprintf(playerConn, "%v", err)
		return nil
	}

	if len(arguments) == 0 {

		fmt.Fprintf(playerConn, "%s\n", currentRoom.Area.Name)
		fmt.Fprintf(playerConn, "%s\n", currentRoom.Area.Description)
		fmt.Fprintf(playerConn, "-----------------------\n\n")
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
	return nil
}

func HandleTakeCommand(db *sql.DB, player interfaces.PlayerInterface, command string, arguments []string, currentChannel chan interfaces.ActionInterface, updateChannel func(string)) interface{} {
	roomUUID := player.GetRoom()
	playerConn := player.GetConn()
	currentRoom, err := getRoom(roomUUID, db)
	if err != nil {
		fmt.Fprintf(playerConn, "%v", err)
		return nil
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
	return nil
}

func HandleDropCommand(db *sql.DB, player interfaces.PlayerInterface, command string, arguments []string, currentChannel chan interfaces.ActionInterface, updateChannel func(string)) interface{} {
	roomUUID := player.GetRoom()
	playerConn := player.GetConn()

	playerItems, err := items.GetItemsForPlayer(db, player.GetUUID())
	if err != nil {
		fmt.Fprintf(playerConn, "%v", err)
		return nil
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
	return nil
}

func HandleInventoryCommand(db *sql.DB, player interfaces.PlayerInterface, command string, arguments []string, currentChannel chan interfaces.ActionInterface, updateChannel func(string)) interface{} {
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
	return nil
}

func HandleFooCommand(db *sql.DB, player interfaces.PlayerInterface, command string, arguments []string, currentChannel chan interfaces.ActionInterface, updateChannel func(string)) interface{} {
	return func() {
		fmt.Println("Handling 'foo' command: ", player.GetName())
	}
}

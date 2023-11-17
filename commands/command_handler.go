package commands

import (
	"database/sql"
	"fmt"
	"mud/areas"
	"mud/player"
	"strings"
)

type CommandParser struct {
	commandName string
	arguments   []string
}

type CommandHandler func(*sql.DB, *player.Player, string, []string)

func NewCommandParser(commandString string) *CommandParser {
	// Split the command string into the command name and arguments.
	commandParts := strings.Split(commandString, " ")
	commandName := commandParts[0]
	arguments := commandParts[1:]

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
	"logout": HandleLogoutCommand,
	"hello":  HandleHelloCommand,
	"look":   HandleLookCommand,
	"exits":  HandleExitsCommand,
	"north":  curryMovePlayerCommand("north"),
	"south":  curryMovePlayerCommand("south"),
	"west":   curryMovePlayerCommand("west"),
	"east":   curryMovePlayerCommand("east"),
	"up":     curryMovePlayerCommand("up"),
	"down":   curryMovePlayerCommand("down"),
}

func curryMovePlayerCommand(direction string) CommandHandler {
	return func(db *sql.DB, player *player.Player, commandName string, arguments []string) {
		HandleMovePlayerCommand(db, player, commandName, append(arguments, direction))
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

	return room, nil
}

func HandleMovePlayerCommand(db *sql.DB, player *player.Player, command string, arguments []string) {
	currentRoom, err := getRoom(player.Room, db)
	if err != nil {
		fmt.Fprintf(player.Conn, "%v", err)
		return
	}
	switch command {
	case "north":
		if currentRoom.Exits.North == nil {
			fmt.Fprintf(player.Conn, "You cannot go that way.\n")

		} else {
			fmt.Fprintf(player.Conn, "=======================\n\n")
			player.SetLocation(db, currentRoom.Exits.North.UUID)
			var lookArgs []string
			HandleLookCommand(db, player, "look", lookArgs)
		}
	case "south":
		if currentRoom.Exits.South == nil {
			fmt.Fprintf(player.Conn, "You cannot go that way.\n")

		} else {
			fmt.Fprintf(player.Conn, "=======================\n\n")
			player.SetLocation(db, currentRoom.Exits.South.UUID)
			var lookArgs []string
			HandleLookCommand(db, player, "look", lookArgs)
		}
	case "west":
		if currentRoom.Exits.West == nil {
			fmt.Fprintf(player.Conn, "You cannot go that way.\n")

		} else {
			fmt.Fprintf(player.Conn, "=======================\n\n")
			player.SetLocation(db, currentRoom.Exits.West.UUID)
			var lookArgs []string
			HandleLookCommand(db, player, "look", lookArgs)
		}
	case "east":
		if currentRoom.Exits.East == nil {
			fmt.Fprintf(player.Conn, "You cannot go that way.\n")

		} else {
			fmt.Fprintf(player.Conn, "=======================\n\n")
			player.SetLocation(db, currentRoom.Exits.East.UUID)
			var lookArgs []string
			HandleLookCommand(db, player, "look", lookArgs)
		}
	case "up":
		if currentRoom.Exits.Up == nil {
			fmt.Fprintf(player.Conn, "You cannot go that way.\n")

		} else {
			fmt.Fprintf(player.Conn, "=======================\n\n")
			player.SetLocation(db, currentRoom.Exits.Up.UUID)
			var lookArgs []string
			HandleLookCommand(db, player, "look", lookArgs)
		}
	default:
		if currentRoom.Exits.Down == nil {
			fmt.Fprintf(player.Conn, "You cannot go that way.\n")
		} else {
			fmt.Fprintf(player.Conn, "=======================\n\n")
			player.SetLocation(db, currentRoom.Exits.Down.UUID)
			var lookArgs []string
			HandleLookCommand(db, player, "look", lookArgs)
		}
	}
}

func HandleExitsCommand(db *sql.DB, player *player.Player, command string, arguments []string) {
	currentRoom, err := getRoom(player.Room, db)
	if err != nil {
		fmt.Fprintf(player.Conn, "%v", err)
		return
	}
	if currentRoom.Exits.North != nil {
		northExit, err := getRoom(currentRoom.Exits.North.UUID, db)
		if err != nil {
			fmt.Fprintf(player.Conn, "%v", err)
		}
		fmt.Fprintf(player.Conn, "North: %s\n", northExit.Name)
	}
	if currentRoom.Exits.South != nil {
		southExit, err := getRoom(currentRoom.Exits.South.UUID, db)
		if err != nil {
			fmt.Fprintf(player.Conn, "%v", err)
		}
		fmt.Fprintf(player.Conn, "South: %s\n", southExit.Name)
	}
	if currentRoom.Exits.West != nil {
		westExit, err := getRoom(currentRoom.Exits.West.UUID, db)
		if err != nil {
			fmt.Fprintf(player.Conn, "%v", err)
		}
		fmt.Fprintf(player.Conn, "West: %s\n", westExit.Name)
	}
	if currentRoom.Exits.East != nil {
		eastExit, err := getRoom(currentRoom.Exits.East.UUID, db)
		if err != nil {
			fmt.Fprintf(player.Conn, "%v", err)
		}
		fmt.Fprintf(player.Conn, "East: %s\n", eastExit.Name)
	}
	if currentRoom.Exits.Up != nil {
		upExit, err := getRoom(currentRoom.Exits.Up.UUID, db)
		if err != nil {
			fmt.Fprintf(player.Conn, "%v", err)
		}
		fmt.Fprintf(player.Conn, "Up: %s\n", upExit.Name)
	}
	if currentRoom.Exits.Down != nil {
		downExit, err := getRoom(currentRoom.Exits.Down.UUID, db)
		if err != nil {
			fmt.Fprintf(player.Conn, "%v", err)
		}
		fmt.Fprintf(player.Conn, "Down: %s\n", downExit.Name)
	}
}

func HandleLogoutCommand(db *sql.DB, player *player.Player, command string, arguments []string) {
	fmt.Fprintf(player.Conn, "Goodbye!\n")
	player.Conn.Close()
}

func HandleHelloCommand(db *sql.DB, player *player.Player, command string, arguments []string) {
	playerName := player.Name
	if playerName == "" {
		playerName = "somebody"
	}

	fmt.Fprintf(player.Conn, "Hello, %s!\n", playerName)
}

func HandleLookCommand(db *sql.DB, player *player.Player, command string, arguments []string) {
	currentRoom, err := getRoom(player.Room, db)
	if err != nil {
		fmt.Fprintf(player.Conn, "%v", err)
		return
	}

	if len(arguments) == 0 {
		fmt.Fprintf(player.Conn, "%s\n", currentRoom.Area.Name)
		fmt.Fprintf(player.Conn, "%s\n", currentRoom.Area.Description)
		fmt.Fprintf(player.Conn, "%s\n", currentRoom.Name)
		fmt.Fprintf(player.Conn, "%s\n", currentRoom.Description)
	} else if len(arguments) == 1 {
		switch arguments[0] {
		case "north":
			if currentRoom.Exits.North != nil {
				northExit, err := getRoom(currentRoom.Exits.North.UUID, db)
				if err != nil {
					fmt.Fprintf(player.Conn, "%v", err)
				}
				fmt.Fprintf(player.Conn, "You look North. You see %s\n", northExit.Name)
			} else {
				fmt.Fprintf(player.Conn, "You don't see anything in that direction\n")
			}
		case "south":
			if currentRoom.Exits.South != nil {
				southExit, err := getRoom(currentRoom.Exits.South.UUID, db)
				if err != nil {
					fmt.Fprintf(player.Conn, "%v", err)
				}
				fmt.Fprintf(player.Conn, "You look South. You see %s\n", southExit.Name)
			} else {
				fmt.Fprintf(player.Conn, "You don't see anything in that direction\n")
			}
		case "east":
			if currentRoom.Exits.East != nil {
				eastExit, err := getRoom(currentRoom.Exits.East.UUID, db)
				if err != nil {
					fmt.Fprintf(player.Conn, "%v", err)
				}
				fmt.Fprintf(player.Conn, "You look East. You see %s\n", eastExit.Name)
			} else {
				fmt.Fprintf(player.Conn, "You don't see anything in that direction\n")
			}
		case "west":
			if currentRoom.Exits.West != nil {
				westExit, err := getRoom(currentRoom.Exits.West.UUID, db)
				if err != nil {
					fmt.Fprintf(player.Conn, "%v", err)
				}
				fmt.Fprintf(player.Conn, "You look West. You see %s\n", westExit.Name)
			} else {
				fmt.Fprintf(player.Conn, "You don't see anything in that direction\n")
			}
		case "up":
			if currentRoom.Exits.Up != nil {
				upExit, err := getRoom(currentRoom.Exits.Up.UUID, db)
				if err != nil {
					fmt.Fprintf(player.Conn, "%v", err)
				}
				fmt.Fprintf(player.Conn, "You look Up. You see %s\n", upExit.Name)
			} else {
				fmt.Fprintf(player.Conn, "You don't see anything in that direction\n")
			}
		case "down":
			if currentRoom.Exits.Down != nil {
				downExit, err := getRoom(currentRoom.Exits.Down.UUID, db)
				if err != nil {
					fmt.Fprintf(player.Conn, "%v", err)
				}
				fmt.Fprintf(player.Conn, "You look Down. You see %s\n", downExit.Name)
			} else {
				fmt.Fprintf(player.Conn, "You don't see anything in that direction\n")
			}
		default:
			fmt.Fprintf(player.Conn, "You don't see that.\n")
		}

	} else {
		fmt.Fprintf(player.Conn, "I don't know how to do that yet.\n")
	}
}

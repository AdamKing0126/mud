package commands

import (
	"database/sql"
	"errors"
	"fmt"
	"mud/areas"
	"mud/player"
	"strings"
)

type CommandParser struct {
	commandName string
	arguments   []string
}

type CommandHandler func(*sql.DB, *player.Player, *areas.Area, string, []string)

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
	return func(db *sql.DB, player *player.Player, area *areas.Area, commandName string, arguments []string) {
		HandleMovePlayerCommand(db, player, area, commandName, append(arguments, direction))
	}
}

func getRoomFromArea(roomUUID string, area *areas.Area, db *sql.DB) (*areas.Room, error) {
	var room *areas.Room
	for i := range area.Rooms {
		if area.Rooms[i].UUID == roomUUID {
			room = &area.Rooms[i]
		}
	}
	if room != nil {
		// if we don't have a room, it might be that this room is outside
		// the current area.  So do a search.
		room_rows, err := db.Query("SELECT UUID, area_uuid, name, description, exit_north, exit_south, exit_east, exit_west, exit_up, exit_down from rooms where uuid=?", roomUUID)
		if err != nil {
			return nil, err
		}

		defer room_rows.Close()
		if !room_rows.Next() {
			return nil, fmt.Errorf("room with UUID %s does not exits", roomUUID)
		}

		room = &areas.Room{
			UUID: roomUUID,
		}

		err = room_rows.Scan(&room.UUID, &room.AreaUUID, &room.Name, &room.Description, &room.ExitNorth, &room.ExitSouth, &room.ExitWest, &room.ExitEast, &room.ExitUp, &room.ExitDown)
		if err != nil {
			return nil, err
		}
		return room, nil
	}
	return nil, errors.New("uh oh")
}

func HandleMovePlayerCommand(db *sql.DB, player *player.Player, area *areas.Area, command string, arguments []string) {
	currentRoom, err := getRoomFromArea(player.Room, area, db)
	if err != nil {
		fmt.Fprintf(player.Conn, "%v", err)
		return
	}
	switch command {
	case "north":
		if currentRoom.ExitNorth == "" {
			fmt.Fprintf(player.Conn, "You cannot go that way.\n")

		} else {
			fmt.Fprintf(player.Conn, "=======================\n\n")
			player.SetLocation(db, currentRoom.ExitNorth)
			var lookArgs []string
			HandleLookCommand(db, player, area, "look", lookArgs)
		}
	case "south":
		if currentRoom.ExitSouth == "" {
			fmt.Fprintf(player.Conn, "You cannot go that way.\n")

		} else {
			fmt.Fprintf(player.Conn, "=======================\n\n")
			player.SetLocation(db, currentRoom.ExitSouth)
			var lookArgs []string
			HandleLookCommand(db, player, area, "look", lookArgs)
		}
	case "west":
		if currentRoom.ExitWest == "" {
			fmt.Fprintf(player.Conn, "You cannot go that way.\n")

		} else {
			fmt.Fprintf(player.Conn, "=======================\n\n")
			player.SetLocation(db, currentRoom.ExitWest)
			var lookArgs []string
			HandleLookCommand(db, player, area, "look", lookArgs)
		}
	case "east":
		if currentRoom.ExitEast == "" {
			fmt.Fprintf(player.Conn, "You cannot go that way.\n")

		} else {
			fmt.Fprintf(player.Conn, "=======================\n\n")
			player.SetLocation(db, currentRoom.ExitEast)
			var lookArgs []string
			HandleLookCommand(db, player, area, "look", lookArgs)
		}
	case "up":
		if currentRoom.ExitUp == "" {
			fmt.Fprintf(player.Conn, "You cannot go that way.\n")

		} else {
			fmt.Fprintf(player.Conn, "=======================\n\n")
			player.SetLocation(db, currentRoom.ExitUp)
			var lookArgs []string
			HandleLookCommand(db, player, area, "look", lookArgs)
		}
	default:
		if currentRoom.ExitDown == "" {
			fmt.Fprintf(player.Conn, "You cannot go that way.\n")
		} else {
			fmt.Fprintf(player.Conn, "=======================\n\n")
			player.SetLocation(db, currentRoom.ExitDown)
			var lookArgs []string
			HandleLookCommand(db, player, area, "look", lookArgs)
		}
	}
}

func HandleExitsCommand(db *sql.DB, player *player.Player, area *areas.Area, command string, arguments []string) {
	currentRoom, err := getRoomFromArea(player.Room, area, db)
	if err != nil {
		fmt.Fprintf(player.Conn, "%v", err)
		return
	}
	if currentRoom.ExitNorth != "" {
		northExit, err := getRoomFromArea(currentRoom.ExitNorth, area, db)
		if err != nil {
			fmt.Fprintf(player.Conn, "%v", err)
		}
		fmt.Fprintf(player.Conn, "North: %s\n", northExit.Name)
	}
	if currentRoom.ExitSouth != "" {
		southExit, err := getRoomFromArea(currentRoom.ExitSouth, area, db)
		if err != nil {
			fmt.Fprintf(player.Conn, "%v", err)
		}
		fmt.Fprintf(player.Conn, "South: %s\n", southExit.Name)
	}
	if currentRoom.ExitWest != "" {
		westExit, err := getRoomFromArea(currentRoom.ExitWest, area, db)
		if err != nil {
			fmt.Fprintf(player.Conn, "%v", err)
		}
		fmt.Fprintf(player.Conn, "West: %s\n", westExit.Name)
	}
	if currentRoom.ExitEast != "" {
		eastExit, err := getRoomFromArea(currentRoom.ExitEast, area, db)
		if err != nil {
			fmt.Fprintf(player.Conn, "%v", err)
		}
		fmt.Fprintf(player.Conn, "East: %s\n", eastExit.Name)
	}
	if currentRoom.ExitUp != "" {
		upExit, err := getRoomFromArea(currentRoom.ExitUp, area, db)
		if err != nil {
			fmt.Fprintf(player.Conn, "%v", err)
		}
		fmt.Fprintf(player.Conn, "Up: %s\n", upExit.Name)
	}
	if currentRoom.ExitDown != "" {
		downExit, err := getRoomFromArea(currentRoom.ExitDown, area, db)
		if err != nil {
			fmt.Fprintf(player.Conn, "%v", err)
		}
		fmt.Fprintf(player.Conn, "Down: %s\n", downExit.Name)
	}
}

func HandleLogoutCommand(db *sql.DB, player *player.Player, area *areas.Area, command string, arguments []string) {
	fmt.Fprintf(player.Conn, "Goodbye!\n")
	player.Conn.Close()
}

func HandleHelloCommand(db *sql.DB, player *player.Player, area *areas.Area, command string, arguments []string) {
	playerName := player.Name
	if playerName == "" {
		playerName = "somebody"
	}

	fmt.Fprintf(player.Conn, "Hello, %s!\n", playerName)
}

func HandleLookCommand(db *sql.DB, player *player.Player, area *areas.Area, command string, arguments []string) {
	currentRoom, err := getRoomFromArea(player.Room, area, db)
	if err != nil {
		fmt.Fprintf(player.Conn, "%v", err)
		return
	}

	if len(arguments) == 0 {
		fmt.Fprintf(player.Conn, "%s\n", area.Name)
		fmt.Fprintf(player.Conn, "%s\n", area.Description)
		fmt.Fprintf(player.Conn, "%s\n", currentRoom.Name)
		fmt.Fprintf(player.Conn, "%s\n", currentRoom.Description)
	} else if len(arguments) == 1 {
		switch arguments[0] {
		case "north":
			if currentRoom.ExitNorth != "" {
				northExit, err := getRoomFromArea(currentRoom.ExitNorth, area, db)
				if err != nil {
					fmt.Fprintf(player.Conn, "%v", err)
				}
				fmt.Fprintf(player.Conn, "You look North. You see %s\n", northExit.Name)
			} else {
				fmt.Fprintf(player.Conn, "You don't see anything in that direction\n")
			}
		case "south":
			if currentRoom.ExitSouth != "" {
				southExit, err := getRoomFromArea(currentRoom.ExitSouth, area, db)
				if err != nil {
					fmt.Fprintf(player.Conn, "%v", err)
				}
				fmt.Fprintf(player.Conn, "You look South. You see %s\n", southExit.Name)
			} else {
				fmt.Fprintf(player.Conn, "You don't see anything in that direction\n")
			}
		case "east":
			if currentRoom.ExitEast != "" {
				eastExit, err := getRoomFromArea(currentRoom.ExitEast, area, db)
				if err != nil {
					fmt.Fprintf(player.Conn, "%v", err)
				}
				fmt.Fprintf(player.Conn, "You look East. You see %s\n", eastExit.Name)
			} else {
				fmt.Fprintf(player.Conn, "You don't see anything in that direction\n")
			}
		case "west":
			if currentRoom.ExitWest != "" {
				westExit, err := getRoomFromArea(currentRoom.ExitWest, area, db)
				if err != nil {
					fmt.Fprintf(player.Conn, "%v", err)
				}
				fmt.Fprintf(player.Conn, "You look West. You see %s\n", westExit.Name)
			} else {
				fmt.Fprintf(player.Conn, "You don't see anything in that direction\n")
			}
		case "up":
			if currentRoom.ExitUp != "" {
				upExit, err := getRoomFromArea(currentRoom.ExitUp, area, db)
				if err != nil {
					fmt.Fprintf(player.Conn, "%v", err)
				}
				fmt.Fprintf(player.Conn, "You look Up. You see %s\n", upExit.Name)
			} else {
				fmt.Fprintf(player.Conn, "You don't see anything in that direction\n")
			}
		case "down":
			if currentRoom.ExitDown != "" {
				downExit, err := getRoomFromArea(currentRoom.ExitDown, area, db)
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

package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"mud/areas"
	"mud/commands"
	"mud/display"
	"mud/interfaces"
	"mud/notifications"
	"mud/players"
	"net"
	"strings"
	"sync"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

type CommandRouterInterface interface {
	HandleCommand(db *sql.DB, player interfaces.PlayerInterface, command []byte, currentChannel chan interfaces.ActionInterface, updateChannel func(string))
}

type Server struct {
	connections map[string]interfaces.PlayerInterface
}

func NewServer() *Server {
	return &Server{
		connections: make(map[string]interfaces.PlayerInterface),
	}
}

func (s *Server) handleConnection(conn net.Conn, router CommandRouterInterface, db *sql.DB, areaChannels map[string]chan interfaces.ActionInterface) {
	defer conn.Close()

	fmt.Fprintf(conn, "Welcome! Please enter your player name: ")
	playerName := getPlayerInput(conn)

	player, err := getPlayerFromDB(db, playerName)
	if err != nil {
		fmt.Fprintf(conn, "Error retrieving player info: %v\n", err)
		return
	}

	defer func() {
		err := setPlayerLoggedInStatus(db, player.UUID, false)
		if err != nil {
			fmt.Fprintf(conn, "Error updating player logged_in status: %v\n", err)
		}
	}()

	fmt.Fprintf(conn, "Please enter your password: ")
	passwd := getPlayerInput(conn)
	err = bcrypt.CompareHashAndPassword([]byte(player.GetHashedPassword()), []byte(passwd))
	if err != nil {
		fmt.Fprintf(conn, "Incorrect password.\n")
		return
	}

	player.Conn = conn

	err = setPlayerLoggedInStatus(db, player.UUID, true)
	if err != nil {
		fmt.Fprintf(conn, "Error updating player logged_in status: %v\n", err)
		return
	}
	s.connections[player.UUID] = player
	defer delete(s.connections, player.UUID)

	notifyPlayersInRoomThatNewPlayerHasJoined(player, s.connections)

	if player.Area == "" || player.Room == "" {
		player.Area = "d71e8cf1-d5ba-426c-8915-4c7f5b22e3a9"
		player.Room = "189a729d-4e40-4184-a732-e2c45c66ff46"
	}

	ch := areaChannels[player.Area]

	updateChannel := func(newArea string) {
		ch = areaChannels[newArea]
	}

	router.HandleCommand(db, player, bytes.NewBufferString("look").Bytes(), ch, updateChannel)

	for {
		display.PrintWithColor(player, fmt.Sprintf("\nHP: %d Mana: %d Mvt: %d> ", player.GetHealth(), player.GetMana(), player.GetMovement()), "primary")
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			break
		}

		router.HandleCommand(db, player, buf[:n], ch, updateChannel)
	}
}

func getPlayerInput(reader io.Reader) string {
	r := bufio.NewReader(reader)
	input, _ := r.ReadString('\n')
	return strings.TrimSpace(input)
}

func getPlayerFromDB(db *sql.DB, playerName string) (*players.Player, error) {
	var player players.Player
	var colorProfile = &players.ColorProfile{}
	query := `SELECT p.name, p.uuid, p.area, p.room, p.health, p.health_max, p.movement, p.movement_max, p.mana, p.mana_max, p.password, cp.uuid, cp.name, cp.primary_color, cp.secondary_color, cp.warning_color, cp.danger_color, cp.title_color, cp.description_color
				FROM players p JOIN color_profiles cp ON cp.uuid = p.color_profile
				WHERE p.name = ?`
	err := db.QueryRow(query, playerName).
		Scan(&player.Name, &player.UUID, &player.Area, &player.Room, &player.Health, &player.HealthMax, &player.Movement, &player.MovementMax, &player.Mana, &player.ManaMax, &player.Password, &colorProfile.UUID, &colorProfile.Name, &colorProfile.Primary, &colorProfile.Secondary, &colorProfile.Warning, &colorProfile.Danger, &colorProfile.Title, &colorProfile.Description)
	if err != nil {
		return &player, err
	}

	player.ColorProfile = colorProfile

	return &player, nil
}

func setPlayerLoggedInStatus(db *sql.DB, playerUUID string, loggedIn bool) error {
	_, err := db.Exec("UPDATE players SET logged_in = ? WHERE uuid = ?", loggedIn, playerUUID)
	if err != nil {
		return err
	}
	return nil
}

func notifyPlayersInRoomThatNewPlayerHasJoined(player interfaces.PlayerInterface, connections map[string]interfaces.PlayerInterface) {
	var playersInRoom []interfaces.PlayerInterface
	for _, p := range connections {
		if p.GetRoom() == player.GetRoom() && p.GetUUID() != player.GetUUID() {
			playersInRoom = append(playersInRoom, p)
		}
	}

	for _, p := range playersInRoom {
		fmt.Fprintf(p.GetConn(), "\n%s has entered the room.\n", player.GetName())
		display.PrintWithColor(p, fmt.Sprintf("\nHP: %d Mana: %d Mvt: %d> ", player.GetHealth(), player.GetMana(), player.GetMovement()), "primary")
	}
}

func openDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./sql_database/mud.db")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Println("Database opened successfully")

	return db, nil
}

func main() {
	db, err := openDatabase()
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	router := commands.NewCommandRouter()
	server := NewServer()
	notifier := notifications.NewNotifier(server.connections)

	commands.RegisterCommands(router, notifier, commands.CommandHandlers)
	if len(router.Handlers) == 0 {
		fmt.Println("Warning: no commands registered. Exiting...")
		return
	}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()
	wg := sync.WaitGroup{}

	areaChannels := make(map[string]chan interfaces.ActionInterface)
	rows, err := db.Query("SELECT uuid, name, description FROM areas")
	if err != nil {
		fmt.Println(err)
		return
	}

	areaInstances := make(map[string]interfaces.AreaInterface)
	for rows.Next() {
		var uuid string
		var name string
		var description string
		err := rows.Scan(&uuid, &name, &description)
		if err != nil {
			fmt.Println(err)
			return
		}

		areaInstances[uuid] = areas.NewArea(uuid, name, description)
		areaChannels[uuid] = make(chan interfaces.ActionInterface)
		go areaInstances[uuid].Run(db, areaChannels[uuid], server.connections)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		wg.Add(1)
		go server.handleConnection(conn, router, db, areaChannels)
	}
}

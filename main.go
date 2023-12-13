package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
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

	player := players.NewPlayer(conn)

	defer func() {
		if _, err := db.Exec("UPDATE players SET logged_in = ? WHERE uuid = ?", false, player.UUID); err != nil {
			fmt.Fprintf(conn, "Error updating player logged_in status: %v\n", err)
		}
	}()

	if player.GetName() == "" {
		fmt.Fprintf(conn, "Welcome! Please enter your player name: ")
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		player.Name = strings.TrimSpace(string(buf[:n]))
	}

	var colorProfileUUID string
	err := db.QueryRow("SELECT uuid, area, room, health, color_profile, password FROM players WHERE name = ?", player.Name).
		Scan(&player.UUID, &player.Area, &player.Room, &player.Health, &colorProfileUUID, &player.Password)
	if err != nil {
		fmt.Fprintf(conn, "Error retrieving player info: %v\n", err)
		return
	}

	fmt.Fprintf(conn, "Please enter your password: ")
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	passwd := strings.TrimSpace(string(buf[:n]))
	err = bcrypt.CompareHashAndPassword([]byte(player.GetHashedPassword()), []byte(passwd))
	if err != nil {
		fmt.Fprintf(conn, "Incorrect password.\n")
		return
	}

	_, err = db.Exec("UPDATE players SET logged_in = ? WHERE uuid = ?", true, player.UUID)
	if err != nil {
		fmt.Fprintf(conn, "Error updating player logged_in status: %v\n", err)
		return
	}
	s.connections[player.UUID] = player

	defer delete(s.connections, player.UUID)

	var playersInRoom []interfaces.PlayerInterface
	for _, p := range s.connections {
		if p.GetRoom() == player.GetRoom() && p.GetUUID() != player.GetUUID() {
			playersInRoom = append(playersInRoom, p)
		}
	}

	for _, p := range playersInRoom {
		if otherPlayer, ok := s.connections[p.GetUUID()]; ok {
			fmt.Fprintf(otherPlayer.GetConn(), "\n%s has entered the room.\n", player.GetName())
			display.PrintWithColor(otherPlayer, fmt.Sprintf("\nHP: %d> ", player.GetHealth()), "primary")
		}
	}

	colorProfile, err := players.NewColorProfileFromDB(db, colorProfileUUID)
	if err != nil {
		fmt.Fprintf(conn, "Error retrieving color profile: %v\n", err)
		return
	}

	player.ColorProfile = colorProfile

	if player.Area == "" || player.Room == "" {
		player.Area = "d71e8cf1-d5ba-426c-8915-4c7f5b22e3a9"
		player.Room = "189a729d-4e40-4184-a732-e2c45c66ff46"
	}

	ch := areaChannels[player.Area]

	updateChannel := func(newArea string) {
		ch = areaChannels[newArea]
	}

	router.HandleCommand(db, player, bytes.NewBufferString("look").Bytes(), ch, updateChannel)
	display.PrintWithColor(player, fmt.Sprintf("\nHP: %d> ", player.GetHealth()), "primary")

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			break
		}

		router.HandleCommand(db, player, buf[:n], ch, updateChannel)
		display.PrintWithColor(player, fmt.Sprintf("\nHP: %d> ", player.GetHealth()), "primary")
	}
}

func main() {
	db, err := sql.Open("sqlite3", "./sql_database/mud.db")
	if err != nil {
		log.Fatalf("Failed to open SQLite database: %v", err)
	} else {
		err := db.Ping()
		if err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}
		fmt.Println("Database opened successfully")
	}

	defer db.Close()

	router := commands.NewCommandRouter()
	server := NewServer()
	notifier := notifications.NewNotifier(server.connections)

	commands.RegisterCommands(router, notifier, commands.CommandHandlers)

	// Check if the command router is empty.
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
	rows, err := db.Query("SELECT uuid FROM areas")
	if err != nil {
		fmt.Println(err)
		return
	}

	areaInstances := make(map[string]interfaces.AreaInterface)

	for rows.Next() {
		var uuid string
		err := rows.Scan(&uuid)
		if err != nil {
			fmt.Println(err)
			return
		}

		areaInstances[uuid] = areas.NewArea()
		areaChannels[uuid] = make(chan interfaces.ActionInterface)
		go areaInstances[uuid].Run(db, areaChannels[uuid])
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

	wg.Wait()
}

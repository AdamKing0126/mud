package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"mud/areas"
	"mud/commands"
	"mud/display"
	"mud/interfaces"
	"mud/notifications"
	"mud/players"
	"net"
	"sync"

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

	player, err := players.LoginPlayer(conn, db)
	if err != nil {
		fmt.Fprintf(conn, "Error: %v\n", err)
		return
	}

	defer func() {
		err := player.Logout(db)
		if err != nil {
			fmt.Fprintf(conn, "Error updating player logged_in status: %v\n", err)
		}
	}()

	s.connections[player.GetUUID()] = player
	defer delete(s.connections, player.GetUUID())

	notifyPlayersInRoomThatNewPlayerHasJoined(player, s.connections)

	ch := areaChannels[player.GetArea()]

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

	areaInstances := make(map[string]*areas.Area)
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
		/// is this interfaces.ActionInterface the problem?
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

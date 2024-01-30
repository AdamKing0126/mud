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
	"mud/world_state"
	"net"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type CommandRouterInterface interface {
	HandleCommand(db *sql.DB, player interfaces.Player, command []byte, currentChannel chan interfaces.Action, updateChannel func(string))
}

type Server struct {
	connections map[string]interfaces.Player
}

func NewServer() *Server {
	return &Server{
		connections: make(map[string]interfaces.Player),
	}
}

func (s *Server) handleConnection(conn net.Conn, router CommandRouterInterface, db *sql.DB, areaChannels map[string]chan interfaces.Action, roomToAreaMap map[string]string, worldState *world_state.WorldState) {
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

	currentRoom := worldState.GetRoom(player.GetRoomUUID(), false)
	currentRoom.AddPlayer(player)
	// TODO I removed this, we probably don't need this function inside WorldState, right?
	// worldState.AddPlayerToRoom(currentRoom.GetUUID(), player)
	// No, currently as it stands we still need taht functionality, at least while
	// moving players from one room to another.
	player.Room = currentRoom
	notifyPlayersInRoomThatNewPlayerHasJoined(player, s.connections)

	ch := areaChannels[player.GetAreaUUID()]

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

func notifyPlayersInRoomThatNewPlayerHasJoined(player interfaces.Player, connections map[string]interfaces.Player) {
	var playersInRoom []interfaces.Player
	for _, p := range connections {
		if p.GetRoomUUID() == player.GetRoomUUID() && p.GetUUID() != player.GetUUID() {
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

func loadAreas(db *sql.DB, server *Server) (map[string]interfaces.Area, map[string]string, map[string]chan interfaces.Action, error) {
	areaInstances := make(map[string]*areas.Area)
	areaInstancesInterface := make(map[string]interfaces.Area)
	roomToAreaMap := make(map[string]string)
	areaChannels := make(map[string]chan interfaces.Action)

	queryString := `
		SELECT r.uuid, a.uuid, a.name, a.description 
		FROM rooms r
		JOIN areas a ON r.area_uuid = a.uuid;
	`
	rows, err := db.Query(queryString)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error retrieving areas/rooms: %v", err)
	}
	for rows.Next() {
		var roomUUID, areaUUID, name, description string
		err := rows.Scan(&roomUUID, &areaUUID, &name, &description)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("error scanning areas/rooms: %v", err)
		}
		_, ok := areaInstances[areaUUID]
		if !ok {
			areaInstances[areaUUID] = areas.NewArea(areaUUID, name, description)
			areaInstancesInterface[areaUUID] = areaInstances[areaUUID]
			areaChannels[areaUUID] = make(chan interfaces.Action)
			go areaInstances[areaUUID].Run(db, areaChannels[areaUUID], server.connections)
		}
		roomToAreaMap[roomUUID] = areaUUID
	}

	return areaInstancesInterface, roomToAreaMap, areaChannels, nil
}

func logoutAllPlayers(db *sql.DB) {
	queryString := `
		UPDATE players SET logged_in = 0 WHERE logged_in = 1;
	`
	_, err := db.Exec(queryString)
	if err != nil {
		fmt.Printf("error logging out players: %v", err)
	}
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

	logoutAllPlayers(db)
	areaInstances, roomToAreaMap, areaChannels, err := loadAreas(db, server)
	if err != nil {
		fmt.Printf("error loading areas: %v", err)
	}

	worldState := world_state.NewWorldState(areaInstances, roomToAreaMap, db)

	commands.RegisterCommands(router, notifier, worldState, commands.CommandHandlers)
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

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		wg.Add(1)
		go server.handleConnection(conn, router, db, areaChannels, roomToAreaMap, worldState)
	}
}

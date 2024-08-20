package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/adamking0126/mud/areas"
	"github.com/adamking0126/mud/commands"
	"github.com/adamking0126/mud/display"
	"github.com/adamking0126/mud/notifications"
	"github.com/adamking0126/mud/players"
	worldState "github.com/adamking0126/mud/world_state"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type CommandRouterInterface interface {
	HandleCommand(db *sqlx.DB, player *players.Player, command []byte, currentChannel chan areas.Action, updateChannel func(string))
}

type Server struct {
	connections map[string]*players.Player
}

func NewServer() *Server {
	return &Server{
		connections: make(map[string]*players.Player),
	}
}

func (s *Server) handleConnection(session ssh.Session, router CommandRouterInterface, db *sqlx.DB, areaChannels map[string]chan areas.Action, roomToAreaMap map[string]string, worldState *worldState.WorldState) {
	defer session.Close()

	player, err := players.LoginPlayer(session, db)
	if err != nil {
		fmt.Fprintf(session, "Error: %v\n", err)
		return
	}
	if player == nil {
		return
	}

	defer func() {
		err := player.Logout(db)
		if err != nil {
			fmt.Fprintf(session, "Error updating player logged_in status: %v\n", err)
		}
	}()

	s.connections[player.UUID] = player
	defer delete(s.connections, player.UUID)

	currentRoom := worldState.GetRoom(player.RoomUUID, false)
	currentRoom.AddPlayer(player)

	notifyPlayersInRoomThatNewPlayerHasJoined(player, s.connections)

	ch := areaChannels[player.AreaUUID]

	updateChannel := func(newArea string) {
		ch = areaChannels[newArea]
	}

	router.HandleCommand(db, player, bytes.NewBufferString("look").Bytes(), ch, updateChannel)

	for {
		display.PrintWithColor(player, fmt.Sprintf("\nHP: %d Mvt: %d> ", player.HP, player.Movement), "primary")
		buf := make([]byte, 1024)
		n, err := session.Read(buf)
		if err != nil {
			fmt.Println(err)
			break
		}

		router.HandleCommand(db, player, buf[:n], ch, updateChannel)
	}
}

func notifyPlayersInRoomThatNewPlayerHasJoined(player *players.Player, connections map[string]*players.Player) {
	var playersInRoom []*players.Player
	for _, p := range connections {
		if p.RoomUUID == player.RoomUUID && p.UUID != player.UUID {
			playersInRoom = append(playersInRoom, p)
		}
	}

	for _, p := range playersInRoom {
		fmt.Fprintf(p.GetSession(), "\n%s has joined the game.\n", player.Name)
		display.PrintWithColor(p, fmt.Sprintf("\nHP: %d Mvt: %d> ", player.HP, player.Movement), "primary")
	}
}

func openDatabase() (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", "./sql_database/mud.db")
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

// TODO should this function be moved into the world_state package?
func loadAreas(db *sqlx.DB, server *Server) (map[string]*areas.Area, map[string]string, map[string]chan areas.Action, error) {
	areaInstances := make(map[string]*areas.Area)
	areaInstancesInterface := make(map[string]*areas.Area)
	roomToAreaMap := make(map[string]string)
	areaChannels := make(map[string]chan areas.Action)

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
			areaChannels[areaUUID] = make(chan areas.Action)
			go areaInstances[areaUUID].Run(db, areaChannels[areaUUID], server.connections)
		}
		roomToAreaMap[roomUUID] = areaUUID
	}

	return areaInstancesInterface, roomToAreaMap, areaChannels, nil
}

func logoutAllPlayers(db *sqlx.DB) {
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
		log.Fatalln(err)
	}
	defer db.Close()

	server := NewServer()
	notifier := notifications.NewNotifier(server.connections)

	logoutAllPlayers(db)
	areaInstances, roomToAreaMap, areaChannels, err := loadAreas(db, server)
	if err != nil {
		log.Fatalf("error loading areas: %v", err)
	}

	worldState := worldState.NewWorldState(areaInstances, roomToAreaMap, db)

	s, err := wish.NewServer(
		wish.WithAddress(":2222"),
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithMiddleware(
			BubbleteaMUD(db, server, notifier, areaChannels, roomToAreaMap, worldState),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Starting SSH server on :2222")
	log.Fatalln(s.ListenAndServe())
}

type mudModel struct {
	db            *sqlx.DB
	server        *Server
	notifier      *notifications.Notifier
	areaChannels  map[string]chan areas.Action
	roomToAreaMap map[string]string
	worldState    *worldState.WorldState
	router        CommandRouterInterface
	player        *players.Player
	session       ssh.Session

	// new fields for state management
	currentState gameState
	loginState   loginState
	charState    characterState
	gameState    playState
}

type gameState int

const (
	stateLogin gameState = iota
	stateCharacter
	statePlay
)

type loginState struct {
	username string
	error    string
}

type characterState struct {
	characters  []*players.Player
	newCharName string
	error       string
}

type playState struct {
	currentRoom *areas.Room
	// add other gameplay related fields as needed
}

func (m mudModel) Init() tea.Cmd {
	return nil
}

func (m mudModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.currentState {
	case stateLogin:
		return m.updateLogin(msg)
	case stateCharacter:
		return m.updateCharacter(msg)
	case statePlay:
		return m.updatePlay(msg)
	}

	return m, cmd
}

// Implement these methods next
func (m *mudModel) updateLogin(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: Implement login logic
	return m, nil
}

func (m *mudModel) updateCharacter(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: Implement character selection/creation logic
	return m, nil
}

func (m *mudModel) updatePlay(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: Implement gameplay logic
	return m, nil
}

func (m mudModel) viewLogin() string {
	// TODO: Implement login view
	return "Login View"
}

func (m mudModel) viewCharacter() string {
	// TODO: Implement character selection/creation view
	return "Character View"
}

func (m mudModel) viewPlay() string {
	// TODO: Implement gameplay view
	return "Game View"
}

func (m mudModel) View() string {
	switch m.currentState {
	case stateLogin:
		return m.viewLogin()
	case stateCharacter:
		return m.viewCharacter()
	case statePlay:
		return m.viewPlay()
	default:
		return "Loading..."
	}
}

func BubbleteaMUD(db *sqlx.DB, server *Server, notifier *notifications.Notifier, areaChannels map[string]chan areas.Action, roomToAreaMap map[string]string, worldState *worldState.WorldState) wish.Middleware {
	return func(sh ssh.Handler) ssh.Handler {
		return func(s ssh.Session) {
			player, err := players.LoginPlayer(s, db)
			if err != nil || player == nil {
				fmt.Fprintf(s, "Login failed: %v\n", err)
				return
			}

			router := commands.NewCommandRouter()
			commands.RegisterCommands(router, notifier, worldState, commands.CommandHandlers)

			m := mudModel{
				db:            db,
				server:        server,
				notifier:      notifier,
				router:        router,
				player:        player,
				session:       s,
				areaChannels:  areaChannels,
				roomToAreaMap: roomToAreaMap,
				worldState:    worldState,
			}

			p := tea.NewProgram(m)
			if _, err := p.Run(); err != nil {
				log.Println("Error running program:", err)
			}

			player.Logout(db)
		}
	}
}

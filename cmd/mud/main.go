package main

import (
	"context"
	"fmt"
	"log"

	"github.com/adamking0126/mud/internal/commands"
	"github.com/adamking0126/mud/internal/display"
	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/internal/game/players"
	"github.com/adamking0126/mud/internal/game/world_state"
	worldState "github.com/adamking0126/mud/internal/game/world_state"
	"github.com/adamking0126/mud/internal/notifications"
	"github.com/adamking0126/mud/internal/ui/login"
	"github.com/adamking0126/mud/pkg/database"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	_ "github.com/mattn/go-sqlite3"

	"github.com/charmbracelet/wish/bubbletea"
)

type CommandRouterInterface interface {
	HandleCommand(
		ctx context.Context,
		worldStateService *worldState.Service,
		playerService *players.Service,
		player *players.Player,
		command []byte,
		currentChannel chan areas.Action,
		updateChannel func(string),
	)
}

type Server struct {
	connections map[string]*players.Player
}

func NewServer() *Server {
	return &Server{
		connections: make(map[string]*players.Player),
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

func main() {
	db, err := database.NewSQLiteDB("./pkg/database/sqlite_databases/mud.db")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	ctx := context.Background()
	server := NewServer()
	notifier := notifications.NewNotifier(server.connections)

	playerService := players.NewService(db)
	err = playerService.LogoutAllPlayers(ctx)
	if err != nil {
		log.Fatalf("error logging out all players: %v", err)
	}

	areaService := areas.NewService(db, playerService, server.connections)
	areaService.LoadAreas(ctx)

	worldStateService := world_state.NewService(db, areaService)

	loginModel := login.NewModel()
	s, err := wish.NewServer(
		wish.WithAddress(":2222"),
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithMiddleware(
			MUDMiddleware(ctx, db, server, notifier, worldStateService, playerService, areaService, loginModel),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Starting SSH server on :2222")
	log.Fatalln(s.ListenAndServe())
}

type mudModel struct {
	db            database.DB
	server        *Server
	notifier      *notifications.Notifier
	areaChannels  map[string]chan areas.Action
	roomToAreaMap map[string]string
	router        CommandRouterInterface
	player        *players.Player
	session       ssh.Session

	// new fields for state management
	currentState      gameState
	loginState        loginState
	charState         characterState
	gameState         playState
	playerService     *players.Service
	areaService       *areas.Service
	worldStateService *worldState.Service

	// bubbletea sub-applications
	loginModel *login.Model
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
	return tea.EnterAltScreen // force a render
}

func (m mudModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.player == nil {
		m.currentState = stateLogin
	}

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

func (m *mudModel) updateLogin(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.loginModel == nil {
		loginModel := login.NewModel()
		m.loginModel = loginModel
	}

	model, cmd := m.loginModel.Update(msg)
	if loginModel, ok := model.(*login.Model); ok {
		m.loginModel = loginModel
	}

	if m.loginModel.Player != nil {
		m.player = m.loginModel.Player
		m.currentState = statePlay
		m.loginModel = nil
	}

	// shoudl I return m or should I return model? TODO
	// if i return m, I don't get the user input displayed in the ssh session.
	return model, cmd
}

// Implement these methods next
func (m *mudModel) updateCharacter(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: Implement character selection/creation logic
	return m, nil
}

func (m *mudModel) updatePlay(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: Implement gameplay logic
	return m, nil
}

func (m mudModel) viewLogin() string {
	if m.loginModel != nil {
		return m.loginModel.View()
	}
	return ""
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
	if m.player == nil {
		return m.viewLogin()
	}

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

func MUDMiddleware(ctx context.Context, db database.DB, server *Server, notifier *notifications.Notifier, worldStateService *worldState.Service, playerService *players.Service, areaService *areas.Service, loginModel *login.Model) wish.Middleware {
	return func(sh ssh.Handler) ssh.Handler {
		return func(s ssh.Session) {
			// player, err := players.LoginPlayer(ctx, s, playerService)
			// if err != nil || player == nil {
			// 	fmt.Fprintf(s, "Login failed: %v\n", err)
			// 	return
			// }

			router := commands.NewCommandRouter()
			commands.RegisterCommands(
				router,
				notifier,
				worldStateService,
				playerService,
				areaService,
				commands.CommandHandlers,
			)

			m := mudModel{
				db:                db,
				server:            server,
				notifier:          notifier,
				router:            router,
				player:            nil,
				session:           s,
				worldStateService: worldStateService,
				playerService:     playerService,
				areaService:       areaService,
				loginModel:        loginModel,
			}
			opts := append(bubbletea.MakeOptions(s), tea.WithAltScreen())

			p := tea.NewProgram(m, opts...)
			if _, err := p.Run(); err != nil {
				log.Println("Error running program:", err)
			}

			// playerService.LogoutPlayer(ctx, m.player)
		}
	}

}

// WTF is this thing even used at all?
// func (s *Server) handleConnection(
// 	ctx context.Context,
// 	session ssh.Session,
// 	router CommandRouterInterface,
// 	db database.DB,
// 	areaChannels map[string]chan areas.Action,
// 	roomToAreaMap map[string]string,
// 	worldState *worldState.WorldState,
// ) {
// 	defer session.Close()

// 	player, err := players.LoginPlayer(ctx, session, db)
// 	if err != nil {
// 		fmt.Fprintf(session, "Error: %v\n", err)
// 		return
// 	}
// 	if player == nil {
// 		return
// 	}
// 	defer func() {
// 		err := player.Logout(ctx, db)
// 		if err != nil {
// 			fmt.Fprintf(session, "Error updating player logged_in status: %v\n", err)
// 		}
// 	}()
// 	s.connections[player.UUID] = player
// 	defer delete(s.connections, player.UUID)
// 	currentRoom := worldState.GetRoom(ctx, player.RoomUUID, false)
// 	currentRoom.AddPlayer(player)
// 	notifyPlayersInRoomThatNewPlayerHasJoined(player, s.connections)
// 	ch := areaChannels[player.AreaUUID]
// 	updateChannel := func(newArea string) {
// 		ch = areaChannels[newArea]
// 	}
// 	router.HandleCommand(ctx, db, player, bytes.NewBufferString("look").Bytes(), ch, updateChannel)
// 	for {
// 		display.PrintWithColor(player, fmt.Sprintf("\nHP: %d Mvt: %d> ", player.HP, player.Movement), "primary")
// 		buf := make([]byte, 1024)
// 		n, err := session.Read(buf)
// 		if err != nil {
// 			fmt.Println(err)
// 			break
// 		}
// 		router.HandleCommand(ctx, db, player, buf[:n], ch, updateChannel)
// 	}
// }

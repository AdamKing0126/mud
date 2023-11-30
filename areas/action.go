package areas

import (
	"database/sql"
	"fmt"
	"mud/interfaces"
	"time"
)

type Action struct {
	Player    interfaces.PlayerInterface
	Command   string
	Arguments []string
}

func (a *Action) GetPlayer() interfaces.PlayerInterface {
	return a.Player
}

func (a *Action) GetCommand() string {
	return a.Command
}

func (a *Action) GetArguments() []string {
	return a.Arguments
}

type ActionHandler interface {
	Execute(db *sql.DB, player interfaces.PlayerInterface, action *Action, updateChannel func(string))
}

type FooActionHandler struct{}

func (h *FooActionHandler) Execute(db *sql.DB, player interfaces.PlayerInterface, action *Action, updateChannel func(string)) {
	playerConn := player.GetConn()
	fmt.Fprintf(playerConn, "FooActionHandler.Execute()\n")
}

// type ActionHandler func(db *sql.DB, action interfaces.ActionInterface)

var ActionHandlers = map[string]ActionHandler{
	"foo": &FooActionHandler{},
}

func (a *Area) Run(db *sql.DB, ch chan interfaces.ActionInterface) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	playerActions := make(map[interfaces.PlayerInterface][]interfaces.ActionInterface)

	for {
		select {
		case action := <-ch:
			player := action.GetPlayer()
			playerActions[player] = append(playerActions[player], action)
		case <-ticker.C:
			// Process one action for each player
			for player, actions := range playerActions {
				if len(actions) > 0 {
					action := actions[0]
					playerActions[player] = actions[1:]

					actionConcrete, ok := action.(*Action)
					if !ok {
						fmt.Println("Unknown action type: ", action)
						continue

					}

					handler, ok := ActionHandlers[actionConcrete.GetCommand()]
					if !ok {
						fmt.Println("Unknown action command: ", actionConcrete.GetCommand())
						continue

					}
					handler.Execute(db, player, actionConcrete, func(message string) {
						fmt.Println("Running command: ", action.GetCommand())
					})
				} else {
					fmt.Println("No commands to run for player.")
				}
			}
		}
	}
}

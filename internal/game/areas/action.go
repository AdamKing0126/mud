package areas

import (
	"fmt"
	"time"

	"github.com/adamking0126/mud/internal/display"
	"github.com/adamking0126/mud/internal/game/players"

	"github.com/jmoiron/sqlx"
)

type Action struct {
	Player    players.Player
	Command   string
	Arguments []string
}

func (a *Action) GetPlayer() players.Player {
	return a.Player
}

func (a *Action) GetCommand() string {
	return a.Command
}

func (a *Action) GetArguments() []string {
	return a.Arguments
}

type ActionHandler interface {
	Execute(db *sqlx.DB, player players.Player, action Action, updateChannel func(string))
}

// TODO WTF is this?
// type FooActionHandler struct{}

// func (h *FooActionHandler) Execute(db *sqlx.DB, player players.Player, action *Action, updateChannel func(string)) {
// 	display.PrintWithColor(player, "FooActionHandler.Execute()\n", "danger")
// }

// var ActionHandlers = map[string]ActionHandler{
// 	"foo": &FooActionHandler{},
// }

var ActionHandlers = map[string]ActionHandler{}

type PlayerActions struct {
	Player  *players.Player
	Actions []Action
}

func (a *Area) Run(db *sqlx.DB, ch chan Action, connections map[string]*players.Player) {
	ticker := time.NewTicker(time.Second)
	tickerCounter := 0
	defer ticker.Stop()

	playerActionsMap := make(map[string]PlayerActions)

	for {
		select {
		case action := <-ch:
			player := action.GetPlayer()
			pa := playerActionsMap[player.UUID]
			pa.Actions = append(pa.Actions, action)
			playerActionsMap[player.UUID] = pa
		case <-ticker.C:
			tickerCounter++
			if tickerCounter%15 == 0 {
				var playersInArea []*players.Player
				playersInArea = make([]*players.Player, 0, len(connections))
				for _, player := range connections {
					areaUUID := a.UUID
					playerArea := player.AreaUUID
					if playerArea == areaUUID {
						playersInArea = append(playersInArea, player)
					}
				}
				for _, player := range playersInArea {
					// Process what hapens on the beat.
					display.PrintWithColor(player, "\nboom-boom\n", "danger")
					if err := player.Regen(db); err != nil {
						fmt.Printf("Error: %v\n", err)
					}
					display.PrintWithColor(player, fmt.Sprintf("\nHP: %d Mvt: %d> ", player.HP, player.Movement), "primary")
				}
			}

			// Process one action for each player
			for _, playerActions := range playerActionsMap {
				if len(playerActions.Actions) > 0 {
					action := playerActions.Actions[0]
					playerActions.Actions = playerActions.Actions[1:]

					handler, ok := ActionHandlers[action.GetCommand()]
					if !ok {
						fmt.Println("Unknown action command: ", action.GetCommand())
						continue

					}
					handler.Execute(db, *playerActions.Player, action, func(message string) {
						fmt.Println("Running command: ", action.GetCommand())
					})
				} else {
					fmt.Println("No commands to run for player.")
				}
			}
		}
	}
}

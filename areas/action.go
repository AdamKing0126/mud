package areas

import (
	"fmt"
	"mud/display"
	"mud/interfaces"
	"time"

	"github.com/jmoiron/sqlx"
)

type Action struct {
	Player    interfaces.Player
	Command   string
	Arguments []string
}

func (a *Action) GetPlayer() interfaces.Player {
	return a.Player
}

func (a *Action) GetCommand() string {
	return a.Command
}

func (a *Action) GetArguments() []string {
	return a.Arguments
}

type ActionHandler interface {
	Execute(db *sqlx.DB, player interfaces.Player, action *Action, updateChannel func(string))
}

type FooActionHandler struct{}

func (h *FooActionHandler) Execute(db *sqlx.DB, player interfaces.Player, action *Action, updateChannel func(string)) {
	display.PrintWithColor(player, "FooActionHandler.Execute()\n", "danger")
}

var ActionHandlers = map[string]ActionHandler{
	"foo": &FooActionHandler{},
}

func (a *Area) Run(db *sqlx.DB, ch chan interfaces.Action, connections map[string]interfaces.Player) {
	ticker := time.NewTicker(time.Second)
	tickerCounter := 0
	defer ticker.Stop()

	playerActions := make(map[interfaces.Player][]interfaces.Action)

	for {
		select {
		case action := <-ch:
			player := action.GetPlayer()
			playerActions[player] = append(playerActions[player], action)
		case <-ticker.C:
			tickerCounter++
			if tickerCounter%15 == 0 {
				var playersInArea []interfaces.Player
				playersInArea = make([]interfaces.Player, 0, len(connections))
				for _, player := range connections {
					areaUUID := a.GetUUID()
					playerArea := player.GetAreaUUID()
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
					display.PrintWithColor(player, fmt.Sprintf("\nHP: %d Mana: %d Mvt: %d> ", player.GetHealth(), player.GetMana(), player.GetMovement()), "primary")
				}
			}

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

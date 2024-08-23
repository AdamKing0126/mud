package commands

import (
	"context"

	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/internal/game/players"
	worldState "github.com/adamking0126/mud/internal/game/world_state"
	"github.com/adamking0126/mud/internal/notifications"
)

type CommandHandler interface {
	Execute(
		ctx context.Context,
		player *players.Player,
		command string,
		arguments []string,
		currentChannel chan areas.Action,
		updateChannel func(string),
	)
}

type CommandHandlerWithPriority struct {
	Handler           CommandHandler
	Notifier          notifications.Notifier
	WorldStateService worldState.Service
	Priority          int
}

var CommandHandlers = map[string]CommandHandlerWithPriority{
	"north":      {Handler: &MovePlayerCommandHandler{Direction: "north"}, Priority: 1},
	"south":      {Handler: &MovePlayerCommandHandler{Direction: "south"}, Priority: 1},
	"west":       {Handler: &MovePlayerCommandHandler{Direction: "west"}, Priority: 1},
	"east":       {Handler: &MovePlayerCommandHandler{Direction: "east"}, Priority: 1},
	"up":         {Handler: &MovePlayerCommandHandler{Direction: "up"}, Priority: 1},
	"down":       {Handler: &MovePlayerCommandHandler{Direction: "down"}, Priority: 1},
	"say":        {Handler: &SayHandler{}, Priority: 2},
	"'":          {Handler: &SayHandler{}, Priority: 2},
	"tell":       {Handler: &TellHandler{}, Priority: 2},
	"give":       {Handler: &GiveCommandHandler{}, Priority: 2},
	"look":       {Handler: &LookCommandHandler{}, Priority: 2},
	"area":       {Handler: &AreaCommandHandler{}, Priority: 2},
	"logout":     {Handler: &LogoutCommandHandler{}, Priority: 10},
	"exits":      {Handler: &ExitsCommandHandler{}, Priority: 2},
	"take":       {Handler: &TakeCommandHandler{}, Priority: 3},
	"drop":       {Handler: &DropCommandHandler{}, Priority: 2},
	"inventory":  {Handler: &InventoryCommandHandler{}, Priority: 2},
	"foo":        {Handler: &FooCommandHandler{}, Priority: 2},
	"/sethealth": {Handler: &AdminSetHealthCommandHandler{}, Priority: 10},
	"status":     {Handler: &PlayerStatusCommandHandler{}, Priority: 2},
	"equip":      {Handler: &EquipHandler{}, Priority: 2},
	"remove":     {Handler: &RemoveCommandHandler{}, Priority: 2},
	"whoami":     {Handler: &WhoAmICommandHandler{}, Priority: 10},
}

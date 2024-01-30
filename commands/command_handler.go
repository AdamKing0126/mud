package commands

import (
	"database/sql"
	"fmt"
	"mud/areas"
	"mud/combat"
	"mud/display"
	"mud/interfaces"
	"mud/items"
	"mud/notifications"
	"mud/players"
	"mud/utils"
	"mud/world_state"
	"reflect"
	"strconv"
	"strings"
)

var CommandHandlers = map[string]utils.CommandHandlerWithPriority{
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
	"status":     {Handler: &PlayerStatusHandler{}, Priority: 2},
	"equip":      {Handler: &EquipHandler{}, Priority: 2},
	"remove":     {Handler: &RemoveHandler{}, Priority: 2},
	"whoami":     {Handler: &WhoAmIHandler{}, Priority: 10},
}

type WhoAmIHandler struct{}

func (*WhoAmIHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	display.PrintWithColor(player, player.GetName(), "reset")

}

type MovePlayerCommandHandler struct {
	Direction  string
	Notifier   *notifications.Notifier
	WorldState *world_state.WorldState
}

func movePlayerToDirection(worldState *world_state.WorldState, db *sql.DB, player interfaces.Player, room interfaces.Room, direction string, notifier *notifications.Notifier, world_state *world_state.WorldState, currentChannel chan interfaces.Action, updateChannel func(string)) {
	if room == nil || room.GetUUID() == "" {
		display.PrintWithColor(player, "You cannot go that way.", "reset")
	} else {
		display.PrintWithColor(player, "=======================\n\n", "secondary")
		notifier.NotifyRoom(player.GetRoomUUID(), player.GetUUID(), fmt.Sprintf("\n%s goes %s.\n", player.GetName(), direction))

		// TODO is it possible to not have to reference worldState directly?
		// does worldState really need to have any functions like this at all?
		// if we are looking at pointers and can directly modify the data, we shouldn't need to do this, right?
		// currentRoom := player.GetRoom()
		// if err := currentRoom.RemovePlayer(player); err != nil {
		// 	fmt.Printf("Error removing player %s from room %s", player.GetUUID(), currentRoom.GetUUID())
		// }
		// room.AddPlayer(player)
		worldState.RemovePlayerFromRoom(player.GetRoomUUID(), player)
		worldState.AddPlayerToRoom(room.GetUUID(), player)

		notifier.NotifyRoom(room.GetUUID(), player.GetUUID(), fmt.Sprintf("\n%s has arrived.\n", player.GetName()))

		player.SetLocation(db, room.GetUUID())
		var lookArgs []string
		lookHandler := &LookCommandHandler{WorldState: world_state}
		lookHandler.Execute(db, player, "look", lookArgs, currentChannel, updateChannel)
	}
}

func (h *MovePlayerCommandHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	playerRoomUUID := player.GetRoomUUID()
	areaUUID := player.GetAreaUUID()

	currentRoom := h.WorldState.GetRoom(playerRoomUUID, true)
	exits := currentRoom.GetExits()

	switch h.Direction {
	case "north":
		movePlayerToDirection(h.WorldState, db, player, exits.GetNorth(), h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	case "south":
		movePlayerToDirection(h.WorldState, db, player, exits.GetSouth(), h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	case "west":
		movePlayerToDirection(h.WorldState, db, player, exits.GetWest(), h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	case "east":
		movePlayerToDirection(h.WorldState, db, player, exits.GetEast(), h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	case "up":
		movePlayerToDirection(h.WorldState, db, player, exits.GetUp(), h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	default:
		movePlayerToDirection(h.WorldState, db, player, exits.GetDown(), h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	}

	if areaUUID != player.GetAreaUUID() {
		updateChannel(player.GetAreaUUID())
	}
}

func (h *MovePlayerCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

func (h *MovePlayerCommandHandler) SetWorldState(world_state *world_state.WorldState) {
	h.WorldState = world_state
}

type ExitsCommandHandler struct {
	ShowOnlyDirections bool
	WorldState         *world_state.WorldState
}

func (h *ExitsCommandHandler) SetWorldState(world_state *world_state.WorldState) {
	h.WorldState = world_state
}

func (h *ExitsCommandHandler) Execute(_ *sql.DB, player interfaces.Player, _ string, _ []string, _ chan interfaces.Action, _ func(string)) {
	roomUUID := player.GetRoomUUID()
	currentRoom := h.WorldState.GetRoom(roomUUID, true)
	exits := currentRoom.GetExits()
	exitMap := map[string]interfaces.Room{
		"North": exits.GetNorth(),
		"South": exits.GetSouth(),
		"West":  exits.GetWest(),
		"East":  exits.GetEast(),
		"Up":    exits.GetUp(),
		"Down":  exits.GetDown(),
	}

	abbreviatedDirections := []string{}
	longDirections := []string{}

	for direction, exit := range exitMap {
		if exit != nil {
			abbreviatedDirections = append(abbreviatedDirections, direction)
			exitRoom := h.WorldState.GetRoom(exit.GetUUID(), false)
			longDirections = append(longDirections, fmt.Sprintf("%s: %s", direction, exitRoom.GetName()))
		}
	}
	if h.ShowOnlyDirections {
		display.PrintWithColor(player, fmt.Sprintf("\nExits: %s\n", strings.Join(abbreviatedDirections, ", ")), "reset")
	} else {
		for _, direction := range longDirections {
			display.PrintWithColor(player, fmt.Sprintf("%s\n", direction), "reset")
		}
	}
}

type LogoutCommandHandler struct {
	Notifier   *notifications.Notifier
	WorldState *world_state.WorldState
}

func (h *LogoutCommandHandler) Execute(db *sql.DB, player interfaces.Player, _ string, _ []string, _ chan interfaces.Action, _ func(string)) {
	display.PrintWithColor(player, "Goodbye!\n", "reset")
	if err := player.Logout(db); err != nil {
		fmt.Printf("Error logging out player: %v\n", err)
		return
	}

	err := h.WorldState.RemovePlayerFromRoom(player.GetRoomUUID(), player)
	if err != nil {
		fmt.Printf("error removing player %s from room %s - %v", player.GetUUID(), player.GetRoomUUID(), err)
		return
	}

	h.Notifier.NotifyRoom(player.GetRoomUUID(), player.GetUUID(), fmt.Sprintf("\n%s has left the game.\n", player.GetName()))
}

func (h *LogoutCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

func (h *LogoutCommandHandler) SetWorldState(worldState *world_state.WorldState) {
	h.WorldState = worldState
}

type GiveCommandHandler struct {
	Notifier   *notifications.Notifier
	WorldState *world_state.WorldState
}

func (h *GiveCommandHandler) SetWorldState(worldState *world_state.WorldState) {
	h.WorldState = worldState
}

func (h *GiveCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

func (h *GiveCommandHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	item := player.GetItemFromInventory(arguments[0])
	if item == nil {
		display.PrintWithColor(player, "You don't have that item", "reset")
		return
	}

	currentRoom := player.GetRoom()
	recipient := currentRoom.GetPlayerByName(arguments[1])
	if recipient == nil {
		display.PrintWithColor(player, "You don't see them here", "reset")
		return
	}

	// Todo Adam I built `WorldState#TransferItem` but I don't need to call
	// that function, because I can operate directly on the players.
	// do I need to make some helper function which handles both the Remove/Add?
	// or is this ok?
	player.RemoveItem(item)
	recipient.AddItem(db, item)

	display.PrintWithColor(player, fmt.Sprintf("You give %s to %s\n", item.GetName(), recipient.GetName()), "reset")
	h.Notifier.NotifyPlayer(recipient.GetUUID(), fmt.Sprintf("\n%s gives you %s\n", player.GetName(), item.GetName()))
}

type LookCommandHandler struct {
	WorldState *world_state.WorldState
}

func (h *LookCommandHandler) SetWorldState(world_state *world_state.WorldState) {
	h.WorldState = world_state
}

func (h *LookCommandHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	roomUUID := player.GetRoomUUID()
	currentRoom := h.WorldState.GetRoom(roomUUID, true)
	playerRoom := player.GetRoom()

	if !reflect.DeepEqual(currentRoom, playerRoom) {
		fmt.Printf("whoopsie, currentRoom != playerRoom\n")
	}

	if len(arguments) == 0 {
		display.PrintWithColor(player, fmt.Sprintf("%s\n", currentRoom.GetName()), "primary")
		display.PrintWithColor(player, fmt.Sprintf("%s\n", currentRoom.GetDescription()), "secondary")
		display.PrintWithColor(player, "-----------------------\n\n", "secondary")

		if len(currentRoom.GetItems()) > 0 {
			display.PrintWithColor(player, "You see the following items:\n", "reset")
			for _, item := range currentRoom.GetItems() {
				display.PrintWithColor(player, fmt.Sprintf("%s\n", item.GetName()), "primary")
			}
			display.PrintWithColor(player, "\n", "reset")
		}

		if len(currentRoom.GetPlayers()) > 1 {
			display.PrintWithColor(player, "You see the following players:\n", "reset")
			for _, playerInRoom := range currentRoom.GetPlayers() {
				if player.GetUUID() != playerInRoom.GetUUID() {
					display.PrintWithColor(player, fmt.Sprintf("%s\n", playerInRoom.GetName()), "primary")
				}
			}
			display.PrintWithColor(player, "\n", "reset")
		}

		exitsHandler := &ExitsCommandHandler{ShowOnlyDirections: true, WorldState: h.WorldState}
		exitsHandler.Execute(db, player, "exits", arguments, currentChannel, updateChannel)
	} else if len(arguments) == 1 {
		exits := currentRoom.GetExits()
		exitMap := map[string]interfaces.Room{
			"North": exits.GetNorth(),
			"South": exits.GetSouth(),
			"West":  exits.GetWest(),
			"East":  exits.GetEast(),
			"Up":    exits.GetUp(),
			"Down":  exits.GetDown(),
		}

		lookDirection := arguments[0]
		directionMatch := false

		for direction, exit := range exitMap {
			if lookDirection == direction {
				directionMatch = true
				if exit != nil {
					exitRoom := h.WorldState.GetRoom(exit.GetUUID(), false)
					display.PrintWithColor(player, fmt.Sprintf("You look %s.  You see %s\n", direction, exitRoom.GetName()), "reset")
				} else {
					display.PrintWithColor(player, "You don't see anything in that direction\n", "reset")
				}
			}
		}

		if !directionMatch {
			target := arguments[0]
			found := false
			itemsForPlayer, err := items.GetItemsForPlayer(db, player.GetUUID())
			if err != nil {
				display.PrintWithColor(player, fmt.Sprintf("%v", err), "danger")
			}

			items := append(currentRoom.GetItems(), itemsForPlayer...)
			for _, item := range items {
				if item.GetName() == target {
					display.PrintWithColor(player, fmt.Sprintf("%s\n", item.GetName()), "reset")
					found = true
					break
				}
			}

			for _, playerInRoom := range currentRoom.GetPlayers() {
				if strings.ToLower(playerInRoom.GetName()) == target {
					display.PrintWithColor(player, fmt.Sprintf("You see %s.\n", playerInRoom.GetName()), "reset")
					found = true
					break
				}
			}

			if !found {
				display.PrintWithColor(player, "You don't see that.\n", "reset")
			}
		}
	} else {
		display.PrintWithColor(player, "I don't know how to do that yet.\n", "reset")
	}
}

type AreaCommandHandler struct {
	WorldState *world_state.WorldState
}

func (h *AreaCommandHandler) SetWorldState(world_state *world_state.WorldState) {
	h.WorldState = world_state
}

func (h *AreaCommandHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	area := h.WorldState.GetArea(player.GetAreaUUID())
	display.PrintWithColor(player, fmt.Sprintf("%s\n", area.GetName()), "primary")
	display.PrintWithColor(player, fmt.Sprintf("%s\n", area.GetDescription()), "secondary")
	display.PrintWithColor(player, "-----------------------\n\n", "secondary")
}

type TakeCommandHandler struct {
	Notifier   *notifications.Notifier
	WorldState *world_state.WorldState
}

func (h *TakeCommandHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	roomUUID := player.GetRoomUUID()
	currentRoom := h.WorldState.GetRoom(roomUUID, true)
	items := currentRoom.GetItems()

	if len(items) > 0 {
		for _, item := range items {
			if item.GetName() == arguments[0] {
				if err := currentRoom.RemoveItem(item); err != nil {
					display.PrintWithColor(player, fmt.Sprintf("error removing item from room: %v", err), "reset")
					break
				}
				player.AddItem(db, item)

				display.PrintWithColor(player, fmt.Sprintf("You take the %s.\n", item.GetName()), "reset")
				h.Notifier.NotifyRoom(player.GetRoomUUID(), player.GetUUID(), fmt.Sprintf("\n%s takes %s.\n", player.GetName(), item.GetName()))
				break
			}
		}
	} else {
		display.PrintWithColor(player, "You don't see that here.\n", "reset")
	}
}

func (h *TakeCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

func (h *TakeCommandHandler) SetWorldState(world_state *world_state.WorldState) {
	h.WorldState = world_state
}

type DropCommandHandler struct {
	Notifier   *notifications.Notifier
	WorldState *world_state.WorldState
}

func (h *DropCommandHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	roomUUID := player.GetRoomUUID()
	room := h.WorldState.GetRoom(roomUUID, false)
	playerRoom := player.GetRoom()

	if !reflect.DeepEqual(room, playerRoom) {
		fmt.Printf("whoopsie, currentRoom != playerRoom\n")
	}

	playerItems := player.GetInventory()
	if len(playerItems) > 0 {
		for _, item := range playerItems {
			if item.GetName() == arguments[0] {
				if err := player.RemoveItem(item); err != nil {
					fmt.Printf("error removing item: %s", err)
				}
				room.AddItem(db, item)
				// err := h.WorldState.TransferItem(player, room, item)
				// if err != nil {
				// 	display.PrintWithColor(player, fmt.Sprintf("Failed to update item location: %v\n", err), "danger")
				// }
				query := "UPDATE item_locations SET room_uuid = ?, player_uuid = NULL WHERE item_uuid = ?"
				_, err := db.Exec(query, roomUUID, item.GetUUID())
				if err != nil {
					display.PrintWithColor(player, fmt.Sprintf("Failed to update item location: %v\n", err), "danger")
				}
				display.PrintWithColor(player, fmt.Sprintf("You drop the %s.\n", item.GetName()), "reset")
				h.Notifier.NotifyRoom(player.GetRoomUUID(), player.GetUUID(), fmt.Sprintf("\n%s dropped %s.\n", player.GetName(), item.GetName()))
				break
			}
		}
	} else {
		display.PrintWithColor(player, "You don't have that item.\n", "warning")
	}
}

func (h *DropCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

func (h *DropCommandHandler) SetWorldState(worldState *world_state.WorldState) {
	h.WorldState = worldState
}

type PlayerStatusHandler struct{}

func (h *PlayerStatusHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	playerAbilities := &players.PlayerAbilities{}

	query := "SELECT * FROM player_abilities WHERE player_uuid = ?"
	err := db.QueryRow(query, player.GetUUID()).Scan(&playerAbilities.UUID, &playerAbilities.PlayerUUID, &playerAbilities.Strength, &playerAbilities.Intelligence, &playerAbilities.Wisdom, &playerAbilities.Constitution, &playerAbilities.Charisma, &playerAbilities.Dexterity)
	if err != nil {
		display.PrintWithColor(player, fmt.Sprintf("%v", err), "danger")
	}

	player.SetAbilities(playerAbilities)

	display.PrintWithColor(player, fmt.Sprintf("Strength: %d\n", playerAbilities.GetStrength()), "danger")
	display.PrintWithColor(player, fmt.Sprintf("Dexterity: %d\n", playerAbilities.GetDexterity()), "danger")
	display.PrintWithColor(player, fmt.Sprintf("Constitution: %d\n", playerAbilities.GetConstitution()), "danger")
	display.PrintWithColor(player, fmt.Sprintf("Intelligence: %d\n", playerAbilities.GetIntelligence()), "danger")
	display.PrintWithColor(player, fmt.Sprintf("Wisdom: %d\n", playerAbilities.GetWisdom()), "danger")
	display.PrintWithColor(player, fmt.Sprintf("Charisma: %d\n", playerAbilities.GetCharisma()), "danger")

	// for debugging purposes only - remove later
	display.PrintWithColor(player, "\n\n***********DEBUG***************\n", "danger")
	display.PrintWithColor(player, fmt.Sprintf("Attack Roll Hits: %t\n", combat.AttackRoll(player, player)), "danger")
	display.PrintWithColor(player, "*******************************\n", "danger")
}

type InventoryCommandHandler struct {
	WorldState *world_state.WorldState
}

func (h *InventoryCommandHandler) SetWorldState(world_state *world_state.WorldState) {
	h.WorldState = world_state
}

func (h *InventoryCommandHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	display.PrintWithColor(player, "You are carrying:\n", "secondary")
	playerInventory := player.GetInventory()

	if len(playerInventory) == 0 {
		display.PrintWithColor(player, "Nothing\n", "reset")
	} else {
		for _, item := range playerInventory {
			display.PrintWithColor(player, fmt.Sprintf("%s\n", item.GetName()), "reset")
		}
	}
}

type FooCommandHandler struct {
	Notifier *notifications.Notifier
}

func (h *FooCommandHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	currentChannel <- &areas.Action{Player: player, Command: command, Arguments: arguments}
}

func (h *FooCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

type SayHandler struct {
	Notifier *notifications.Notifier
}

func (h *SayHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	msg := strings.Join(arguments, " ")
	display.PrintWithColor(player, fmt.Sprintf("You say \"%s\"\n", msg), "reset")
	h.Notifier.NotifyRoom(player.GetRoomUUID(), player.GetUUID(), fmt.Sprintf("\n%s says \"%s\"\n", player.GetName(), msg))
}

func (h *SayHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

type TellHandler struct {
	Notifier *notifications.Notifier
}

func (h *TellHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	msg := strings.Join(arguments[1:], " ")
	retrievedPlayer, err := players.GetPlayerByName(db, arguments[0])
	if err != nil {
		display.PrintWithColor(player, fmt.Sprintf("Error retrieving player UUID: %v\n", err), "danger")
		return
	}

	if player.GetUUID() == retrievedPlayer.GetUUID() {
		display.PrintWithColor(player, "Talking to yourself again?\n", "reset")
		return
	}

	if retrievedPlayer.GetLoggedIn() {
		display.PrintWithColor(player, fmt.Sprintf("You tell %s \"%s\"\n", arguments[0], msg), "reset")
		h.Notifier.NotifyPlayer(retrievedPlayer.GetUUID(), fmt.Sprintf("\n%s tells you \"%s\"\n", player.GetName(), msg))
	} else {
		display.PrintWithColor(player, fmt.Sprintf("%s isn't here\n", arguments[0]), "reset")
	}
}

func (h *TellHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

type AdminSetHealthCommandHandler struct {
	Notifier *notifications.Notifier
}

func (h *AdminSetHealthCommandHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	target := arguments[0]
	value := arguments[1]

	retrievedPlayer, err := players.GetPlayerByName(db, target)
	if err != nil {
		display.PrintWithColor(player, fmt.Sprintf("Error retrieving player UUID: %v\n", err), "danger")
		return
	}

	query := "UPDATE players SET health = ? WHERE UUID = ?"

	_, err = db.Exec(query, value, retrievedPlayer.GetUUID())

	if err != nil {
		display.PrintWithColor(player, fmt.Sprintf("Error updating health: %v\n", err), "danger")
		return
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		display.PrintWithColor(player, fmt.Sprintf("Error converting value to int: %v\n", err), "danger")
		return
	}
	h.Notifier.Players[retrievedPlayer.GetUUID()].SetHealth(intValue)
	display.PrintWithColor(player, fmt.Sprintf("You set %s's health to %d\n", target, intValue), "reset")
	h.Notifier.NotifyPlayer(retrievedPlayer.GetUUID(), fmt.Sprintf("\n%s magically sets your health to %d\n", player.GetName(), intValue))

}

func (h *AdminSetHealthCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

type EquipHandler struct {
	Notifier   *notifications.Notifier
	WorldState *world_state.WorldState
}

func (h *EquipHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	if len(arguments) == 0 {
		player.DisplayEquipment()
		return
	}

	playerItems := player.GetInventory()

	if len(playerItems) > 0 {
		for _, item := range playerItems {
			if item.GetName() == arguments[0] {
				if player.Equip(db, item) {
					display.PrintWithColor(player, fmt.Sprintf("You wield %s.\n", item.GetName()), "reset")
					h.Notifier.NotifyRoom(player.GetRoomUUID(), player.GetUUID(), fmt.Sprintf("\n%s wields %s.\n", player.GetName(), item.GetName()))
					player.RemoveItem(item)
					// h.WorldState.RemoveItemFromPlayerInventory(player, item)
				}
				break
			}
		}
	} else {
		display.PrintWithColor(player, "You don't have that item.\n", "warning")
	}

}

func (h *EquipHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

func (h *EquipHandler) SetWorldState(worldState *world_state.WorldState) {
	h.WorldState = worldState
}

type RemoveHandler struct {
	Notifier *notifications.Notifier
}

func (h *RemoveHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

func (h *RemoveHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	if len(arguments) == 0 {
		display.PrintWithColor(player, "Remove what?\n", "primary")
		return
	}

	player.Remove(db, arguments[0])

}

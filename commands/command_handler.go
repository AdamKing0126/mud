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

func movePlayerToDirection(worldState *world_state.WorldState, db *sql.DB, player interfaces.Player, room *areas.Room, direction string, notifier *notifications.Notifier, world_state *world_state.WorldState, currentChannel chan interfaces.Action, updateChannel func(string)) {
	if room == nil || room.UUID == "" {
		display.PrintWithColor(player, "You cannot go that way.", "reset")
	} else {
		display.PrintWithColor(player, "=======================\n\n", "secondary")
		notifier.NotifyRoom(player.GetRoomUUID(), player.GetUUID(), fmt.Sprintf("\n%s goes %s.\n", player.GetName(), direction))
		worldState.RemovePlayerFromRoom(player.GetRoomUUID(), player)

		player.SetLocation(db, room.UUID)
		worldState.AddPlayerToRoom(room.UUID, player)
		notifier.NotifyRoom(room.UUID, player.GetUUID(), fmt.Sprintf("\n%s has arrived.\n", player.GetName()))

		var lookArgs []string
		lookHandler := &LookCommandHandler{WorldState: world_state}
		lookHandler.Execute(db, player, "look", lookArgs, currentChannel, updateChannel)
	}
}

func (h *MovePlayerCommandHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	playerRoomUUID := player.GetRoomUUID()
	areaUUID := player.GetArea()

	currentRoom := h.WorldState.GetRoom(playerRoomUUID, true)

	switch h.Direction {
	case "north":
		movePlayerToDirection(h.WorldState, db, player, currentRoom.Exits.North, h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	case "south":
		movePlayerToDirection(h.WorldState, db, player, currentRoom.Exits.South, h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	case "west":
		movePlayerToDirection(h.WorldState, db, player, currentRoom.Exits.West, h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	case "east":
		movePlayerToDirection(h.WorldState, db, player, currentRoom.Exits.East, h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	case "up":
		movePlayerToDirection(h.WorldState, db, player, currentRoom.Exits.Up, h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	default:
		movePlayerToDirection(h.WorldState, db, player, currentRoom.Exits.Down, h.Direction, h.Notifier, h.WorldState, currentChannel, updateChannel)
	}

	if areaUUID != player.GetArea() {
		updateChannel(player.GetArea())
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

func (h *ExitsCommandHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	roomUUID := player.GetRoomUUID()
	currentRoom := h.WorldState.GetRoom(roomUUID, true)

	exits := map[string]*areas.Room{
		"North": currentRoom.Exits.North,
		"South": currentRoom.Exits.South,
		"West":  currentRoom.Exits.West,
		"East":  currentRoom.Exits.East,
		"Up":    currentRoom.Exits.Up,
		"Down":  currentRoom.Exits.Down,
	}

	abbreviatedDirections := []string{}
	longDirections := []string{}

	for direction, exit := range exits {
		if exit != nil {
			abbreviatedDirections = append(abbreviatedDirections, direction)
			exitRoom := h.WorldState.GetRoom(exit.UUID, false)
			// exitRoom, err := getRoom(exit.UUID, db)
			// if err != nil {
			// 	display.PrintWithColor(player, fmt.Sprintf("%v", err), "danger")
			// }
			longDirections = append(longDirections, fmt.Sprintf("%s: %s", direction, exitRoom.Name))
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
	Notifier *notifications.Notifier
}

func (h *LogoutCommandHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	display.PrintWithColor(player, "Goodbye!\n", "reset")
	if err := player.Logout(db); err != nil {
		fmt.Printf("Error logging out player: %v\n", err)
		return
	}
	h.Notifier.NotifyRoom(player.GetRoomUUID(), player.GetUUID(), fmt.Sprintf("\n%s has left the game.\n", player.GetName()))
}

func (h *LogoutCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
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
	// TODO Adam
	// switch from DB-backed to WorldState-backed.
	// what does that mean, exactly?
	// in WorldState.Areas[areaIdx].Rooms[roomIdx].Players
	// transfer Inventory item from one player to another
	// plus a DB call to save it.
	item, err := items.GetItemByNameForPlayer(db, arguments[0], player.GetUUID())
	if err != nil {
		display.PrintWithColor(player, fmt.Sprintf("%v", err), "danger")
		return
	}
	playersInRoom, err := players.GetPlayersInRoom(db, player.GetRoomUUID())
	if err != nil {
		display.PrintWithColor(player, fmt.Sprintf("%v", err), "danger")
		return
	}

	for _, playerInRoom := range playersInRoom {
		if strings.ToLower(playerInRoom.GetName()) == arguments[1] {
			err := item.SetLocation(db, playerInRoom.GetUUID(), "")
			if err != nil {
				fmt.Println(err)
			}
			display.PrintWithColor(player, fmt.Sprintf("You give %s to %s\n", item.GetName(), arguments[1]), "reset")
			h.Notifier.NotifyPlayer(playerInRoom.GetUUID(), fmt.Sprintf("\n%s gives you %s\n", player.GetName(), item.GetName()))
			return
		}
	}

	display.PrintWithColor(player, fmt.Sprintf("You don't see %s here.\n", arguments[1]), "reset")
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

	if len(arguments) == 0 {
		display.PrintWithColor(player, fmt.Sprintf("%s\n", currentRoom.Name), "primary")
		display.PrintWithColor(player, fmt.Sprintf("%s\n", currentRoom.Description), "secondary")
		display.PrintWithColor(player, "-----------------------\n\n", "secondary")

		if len(currentRoom.Items) > 0 {
			display.PrintWithColor(player, "You see the following items:\n", "reset")
			for _, item := range currentRoom.Items {
				display.PrintWithColor(player, fmt.Sprintf("%s\n", item.GetName()), "primary")
			}
			display.PrintWithColor(player, "\n", "reset")
		}

		if len(currentRoom.Players) > 1 {
			display.PrintWithColor(player, "You see the following players:\n", "reset")
			for _, playerInRoom := range currentRoom.Players {
				if player.GetUUID() != playerInRoom.GetUUID() {
					display.PrintWithColor(player, fmt.Sprintf("%s\n", playerInRoom.GetName()), "primary")
				}
			}
			display.PrintWithColor(player, "\n", "reset")
		}

		exitsHandler := &ExitsCommandHandler{ShowOnlyDirections: true, WorldState: h.WorldState}
		exitsHandler.Execute(db, player, "exits", arguments, currentChannel, updateChannel)
	} else if len(arguments) == 1 {

		exits := map[string]*areas.Room{
			"north": currentRoom.Exits.North,
			"south": currentRoom.Exits.South,
			"west":  currentRoom.Exits.West,
			"east":  currentRoom.Exits.East,
			"down":  currentRoom.Exits.Down,
			"up":    currentRoom.Exits.Up,
		}

		lookDirection := arguments[0]
		directionMatch := false

		for direction, exit := range exits {
			if lookDirection == direction {
				directionMatch = true
				if exit != nil {
					exitRoom := h.WorldState.GetRoom(exit.UUID, false)
					display.PrintWithColor(player, fmt.Sprintf("You look %s.  You see %s\n", direction, exitRoom.Name), "reset")
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

			items := append(currentRoom.Items, itemsForPlayer...)
			for _, item := range items {
				if item.GetName() == target {
					display.PrintWithColor(player, fmt.Sprintf("%s\n", item.GetName()), "reset")
					found = true
					break
				}
			}

			for _, playerInRoom := range currentRoom.Players {
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
	for _, area := range h.WorldState.Areas {
		if area.UUID == player.GetArea() {
			display.PrintWithColor(player, fmt.Sprintf("%s\n", area.Name), "primary")
			display.PrintWithColor(player, fmt.Sprintf("%s\n", area.Description), "secondary")
			display.PrintWithColor(player, "-----------------------\n\n", "secondary")
		}

	}
}

type TakeCommandHandler struct {
	Notifier   *notifications.Notifier
	WorldState *world_state.WorldState
}

func (h *TakeCommandHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	roomUUID := player.GetRoomUUID()
	currentRoom := h.WorldState.GetRoom(roomUUID, true)

	if len(currentRoom.Items) > 0 {
		for _, item := range currentRoom.Items {
			if item.GetName() == arguments[0] {
				err := h.WorldState.TransferItemFromRoomToPlayer(currentRoom, item, player)
				if err != nil {
					fmt.Printf("error transferring item from room to player: %v", err)
				}

				display.PrintWithColor(player, fmt.Sprintf("You take the %s.\n", item.GetName()), "reset")
				h.Notifier.NotifyRoom(player.GetRoomUUID(), player.GetUUID(), fmt.Sprintf("\n%s takes %s.\n", player.GetName(), item.GetName()))
				break
				// item.SetLocation(db, "", player.GetUUID())
				// query := "UPDATE item_locations SET room_uuid = '', player_uuid = ? WHERE item_uuid = ?"
				// _, err := db.Exec(query, player.GetUUID(), item.GetUUID())
				// if err != nil {
				// 	display.PrintWithColor(player, fmt.Sprintf("Failed to update item location: %v\n", err), "danger")
				// }
			}
		}
	} else {
		display.PrintWithColor(player, "You don't see that here.\n", "reset")
	}

	smth := h.WorldState.GetRoom(roomUUID, true)
	fmt.Printf("room: %v", smth.Name)
}

func (h *TakeCommandHandler) SetNotifier(notifier *notifications.Notifier) {
	h.Notifier = notifier
}

func (h *TakeCommandHandler) SetWorldState(world_state *world_state.WorldState) {
	h.WorldState = world_state
}

type DropCommandHandler struct {
	Notifier *notifications.Notifier
}

func (h *DropCommandHandler) Execute(db *sql.DB, player interfaces.Player, command string, arguments []string, currentChannel chan interfaces.Action, updateChannel func(string)) {
	roomUUID := player.GetRoomUUID()

	playerItems, err := items.GetItemsForPlayer(db, player.GetUUID())
	if err != nil {
		display.PrintWithColor(player, fmt.Sprintf("%v", err), "danger")
	}

	if len(playerItems) > 0 {
		for _, item := range playerItems {
			if item.GetName() == arguments[0] {
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
	display.PrintWithColor(player, "You are carryin/g:\n", "secondary")
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
	Notifier *notifications.Notifier
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

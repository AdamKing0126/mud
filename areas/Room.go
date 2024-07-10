package areas

import (
	"fmt"
	"mud/interfaces"
	"mud/items"
	"mud/players"

	"github.com/jmoiron/sqlx"
)

type Room struct {
	UUID        string
	AreaUUID    string
	Name        string
	Description string
	Area        AreaInfo
	Exits       ExitInfo
	Items       []items.Item
	Players     []players.Player
	Mobs        []interfaces.Mob
}

func (room Room) GetUUID() string {
	return room.UUID
}

func (room Room) GetPlayers() []players.Player {
	return room.Players
}

func (room Room) GetMobs() []interfaces.Mob {
	return room.Mobs
}

func (room Room) GetPlayerByName(playerName string) *players.Player {
	playersInRoom := room.GetPlayers()
	for idx := range playersInRoom {
		if playersInRoom[idx].GetName() == playerName {
			return &playersInRoom[idx]
		}
	}
	return nil
}

func (room Room) GetDescription() string {
	return room.Description
}

func (room Room) GetName() string {
	return room.Name
}

func (room Room) GetItems() []items.Item {
	return room.Items
}

func (room Room) AddPlayer(player players.Player) {
	playerIdx := -1
	for idx := range room.Players {
		if room.Players[idx].GetUUID() == player.GetUUID() {
			playerIdx = idx
		}
	}
	if playerIdx == -1 {
		room.Players = append(room.Players, player)
	} else {
		// in case the Players in the room is already loaded from the DB, update with the current player.
		// TODO does this whole logic here need to be reworked so this case doesn't happen?
		room.Players = append(room.Players[:playerIdx], append([]players.Player{player}, room.Players[playerIdx+1:]...)...)
	}
}

func (room Room) RemovePlayer(player players.Player) error {
	playersInRoom := room.GetPlayers()
	for idx, playerInRoom := range playersInRoom {
		if playerInRoom.GetUUID() == player.GetUUID() {
			room.Players = append(room.Players[:idx], room.Players[idx+1:]...)
			return nil
		}
	}
	return fmt.Errorf("player not found")
}

func (room Room) GetExits() ExitInfo {
	return room.Exits
}

func (room *Room) SetExits(exits ExitInfo) {
	room.Exits = exits
}

func (room *Room) SetPlayers(players []players.Player) {
	room.Players = players
}

func (room *Room) SetMobs(mobs []interfaces.Mob) {
	room.Mobs = mobs
}

func (room *Room) SetItems(items []items.Item) {
	room.Items = items
}

func NewAreaInfo(uuid string, name string, description string) *AreaInfo {
	return &AreaInfo{UUID: uuid, Name: name, Description: description}
}

func NewRoomWithAreaInfo(uuid string, area_uuid string, name string, description string, area_name string, area_description string, exit_north string, exit_south string, exit_east string, exit_west string, exit_up string, exit_down string) Room {
	areaInfo := NewAreaInfo(area_uuid, area_name, area_description)
	exitInfo := NewExitInfo(exit_north, exit_south, exit_west, exit_east, exit_up, exit_down)
	return Room{UUID: uuid, AreaUUID: area_uuid, Name: name, Description: description, Area: *areaInfo, Exits: *exitInfo}
}

func (room *Room) AddItem(db *sqlx.DB, item items.Item) error {
	room.Items = append(room.Items, item)
	err := item.SetLocation(db, "", room.UUID)
	if err != nil {
		return err
	}
	return nil
}

func (room *Room) RemoveItem(item items.Item) error {
	items := room.GetItems()
	for itemIndex := range items {
		if items[itemIndex].UUID == item.UUID {
			room.Items = append(room.Items[:itemIndex], room.Items[itemIndex+1:]...)
			return nil
		}
	}
	return fmt.Errorf("item not found")
}

type ExitInfo struct {
	North *Room
	South *Room
	West  *Room
	East  *Room
	Up    *Room
	Down  *Room
}

func (e ExitInfo) GetNorth() *Room {
	return e.North
}

func (e ExitInfo) GetSouth() *Room {
	return e.South
}

func (e ExitInfo) GetWest() *Room {
	return e.West
}

func (e ExitInfo) GetEast() *Room {
	return e.East
}

func (e ExitInfo) GetDown() *Room {
	return e.Down
}

func (e ExitInfo) GetUp() *Room {
	return e.Up
}

func NewExitInfo(north string, south string, west string, east string, up string, down string) *ExitInfo {

	exitInfo := &ExitInfo{}

	if north != "" {
		exitInfo.North = &Room{UUID: north}
	}
	if south != "" {
		exitInfo.South = &Room{UUID: south}
	}
	if west != "" {
		exitInfo.West = &Room{UUID: west}
	}
	if east != "" {
		exitInfo.East = &Room{UUID: east}
	}
	if up != "" {
		exitInfo.Up = &Room{UUID: up}
	}
	if down != "" {
		exitInfo.Down = &Room{UUID: down}
	}
	return exitInfo

}

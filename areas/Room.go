package areas

import (
	"database/sql"
	"fmt"
	"mud/interfaces"
	"mud/mobs"

	"net"
)

type PlayerInRoomInterface interface {
	GetUUID() string
	GetName() string
	GetInventory() []interfaces.Item
	SetInventory([]interfaces.Item)
	GetColorProfile() interfaces.ColorProfile
	GetConn() net.Conn
}

type Room struct {
	UUID        string
	AreaUUID    string
	Name        string
	Description string
	Area        AreaInfo
	Exits       *ExitInfo
	Items       []interfaces.Item
	Players     []interfaces.Player
	Mobs        []mobs.Mob
}

func (room *Room) GetUUID() string {
	return room.UUID
}

func (room *Room) GetPlayers() []interfaces.Player {
	return room.Players
}

func (room *Room) GetPlayerByName(playerName string) interfaces.Player {
	playersInRoom := room.GetPlayers()
	for idx := range playersInRoom {
		if playersInRoom[idx].GetName() == playerName {
			return playersInRoom[idx]
		}
	}
	return nil
}

func (room *Room) GetDescription() string {
	return room.Description
}

func (room *Room) GetName() string {
	return room.Name
}

func (room *Room) GetItems() []interfaces.Item {
	return room.Items
}

func (room *Room) AddPlayer(player interfaces.Player) {
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
		room.Players = append(room.Players[:playerIdx], append([]interfaces.Player{player}, room.Players[playerIdx+1:]...)...)
	}
}

func (room *Room) RemovePlayer(player interfaces.Player) error {
	for idx := range room.GetPlayers() {
		if room.GetPlayers()[idx] == player {
			room.Players = append(room.Players[:idx], room.Players[idx+1:]...)
			return nil
		}
	}
	return fmt.Errorf("player not found")
}

func (room *Room) GetExits() interfaces.ExitInfo {
	return room.Exits
}

func (room *Room) SetExits(exits interfaces.ExitInfo) {
	room.Exits = exits.(*ExitInfo)
}

func (room *Room) SetPlayers(players []interfaces.Player) {
	room.Players = players
}

func (room *Room) SetItems(items []interfaces.Item) {
	room.Items = items
}

func NewAreaInfo(uuid string, name string, description string) *AreaInfo {
	return &AreaInfo{UUID: uuid, Name: name, Description: description}
}

func NewRoomWithAreaInfo(uuid string, area_uuid string, name string, description string, area_name string, area_description string, exit_north string, exit_south string, exit_east string, exit_west string, exit_up string, exit_down string) *Room {
	areaInfo := NewAreaInfo(area_uuid, area_name, area_description)
	exitInfo := NewExitInfo(exit_north, exit_south, exit_west, exit_east, exit_up, exit_down)
	return &Room{UUID: uuid, AreaUUID: area_uuid, Name: name, Description: description, Area: *areaInfo, Exits: exitInfo}
}

func (room *Room) AddItem(db *sql.DB, item interfaces.Item) error {
	room.Items = append(room.Items, item)
	err := item.SetLocation(db, "", room.UUID)
	if err != nil {
		return err
	}
	return nil
}

func (room *Room) RemoveItem(item interfaces.Item) error {
	for itemIndex := range room.GetItems() {
		if room.GetItems()[itemIndex] == item {
			room.Items = append(room.Items[:itemIndex], room.Items[itemIndex+1:]...)
			return nil
		}
	}
	return fmt.Errorf("item not found")
}

type ExitInfo struct {
	North interfaces.Room
	South interfaces.Room
	West  interfaces.Room
	East  interfaces.Room
	Up    interfaces.Room
	Down  interfaces.Room
}

func (e *ExitInfo) GetNorth() interfaces.Room {
	return e.North
}

func (e *ExitInfo) GetSouth() interfaces.Room {
	return e.South
}

func (e *ExitInfo) GetWest() interfaces.Room {
	return e.West
}

func (e *ExitInfo) GetEast() interfaces.Room {
	return e.East
}

func (e *ExitInfo) GetDown() interfaces.Room {
	return e.Down
}

func (e *ExitInfo) GetUp() interfaces.Room {
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

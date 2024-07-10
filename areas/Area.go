package areas

import (
	"fmt"
)

type Area struct {
	UUID        string
	Name        string
	Description string
	Rooms       []Room
	Channel     chan Action
}

func (a Area) GetUUID() string {
	return a.UUID
}

func (a Area) GetName() string {
	return a.Name
}

func (a Area) GetDescription() string {
	return a.Description
}

func (a Area) GetRooms() []Room {
	return a.Rooms
}

func (a Area) GetRoomByUUID(roomUUID string) (*Room, error) {
	rooms := a.GetRooms()
	for idx := range rooms {
		if rooms[idx].GetUUID() == roomUUID {
			return &rooms[idx], nil
		}
	}
	return nil, fmt.Errorf("room UUID %s not found in area %s", roomUUID, a.GetUUID())
}

func (a *Area) SetRooms(rooms []Room) {
	a.Rooms = rooms
}

func (a *Area) SetRoomAtIndex(idx int, room Room) {
	if idx < 0 || idx >= len(a.Rooms) {
		fmt.Printf("Index out of bounds: %d", idx)
		return
	}
	a.Rooms[idx] = room
}

type AreaInfo struct {
	UUID        string
	Name        string
	Description string
}

func NewArea(uuid string, name string, description string) *Area {
	return &Area{UUID: uuid, Name: name, Description: description}
}

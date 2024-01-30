package areas

import (
	"fmt"
	"mud/interfaces"
)

type Area struct {
	UUID        string
	Name        string
	Description string
	Rooms       []interfaces.Room
	Channel     chan Action
}

func (a *Area) GetUUID() string {
	return a.UUID
}

func (a *Area) GetName() string {
	return a.Name
}

func (a *Area) GetDescription() string {
	return a.Description
}

func (a *Area) GetRooms() []interfaces.Room {
	return a.Rooms
}

func (a *Area) GetRoomByUUID(roomUUID string) (interfaces.Room, error) {
	rooms := a.GetRooms()
	for idx := range rooms {
		if rooms[idx].GetUUID() == roomUUID {
			return rooms[idx], nil
		}
	}
	return nil, fmt.Errorf("room UUID %s not found in area %s", roomUUID, a.GetUUID())
}

func (a *Area) SetRooms(rooms []interfaces.Room) {
	a.Rooms = rooms
}

func (a *Area) SetRoomAtIndex(idx int, room interfaces.Room) {
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

package areas

import (
	"fmt"
)

type Area struct {
	UUID        string
	Name        string
	Description string
	Rooms       []*Room
	Channel     chan Action
}

func (a Area) GetRoomByUUID(roomUUID string) (*Room, error) {
	rooms := a.Rooms
	for idx := range rooms {
		if rooms[idx].UUID == roomUUID {
			return rooms[idx], nil
		}
	}
	return nil, fmt.Errorf("room UUID %s not found in area %s", roomUUID, a.UUID)
}

func (a *Area) SetRoomAtIndex(idx int, room *Room) {
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

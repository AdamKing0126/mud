package areas

import (
	"mud/interfaces"
	"mud/mobs"
)

type PlayerInRoomInterface interface {
	GetUUID() string
	GetName() string
}

type Room struct {
	UUID        string
	AreaUUID    string
	Name        string
	Description string
	Area        AreaInfo
	Exits       ExitInfo
	Items       []interfaces.Item
	Players     []PlayerInRoomInterface
	Mobs        []mobs.Mob
}

func NewAreaInfo(uuid string, name string, description string) *AreaInfo {
	return &AreaInfo{UUID: uuid, Name: name, Description: description}
}

func NewRoomWithAreaInfo(uuid string, area_uuid string, name string, description string, area_name string, area_description string, exit_north string, exit_south string, exit_east string, exit_west string, exit_up string, exit_down string) *Room {
	areaInfo := NewAreaInfo(area_uuid, area_name, area_description)
	exitInfo := NewExitInfo(exit_north, exit_south, exit_west, exit_east, exit_up, exit_down)
	return &Room{UUID: uuid, AreaUUID: area_uuid, Name: name, Description: description, Area: *areaInfo, Exits: *exitInfo}
}

type ExitInfo struct {
	North *Room
	South *Room
	West  *Room
	East  *Room
	Up    *Room
	Down  *Room
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

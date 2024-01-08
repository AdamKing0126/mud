package areas

import (
	"mud/interfaces"
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
}

type ExitInfo struct {
	North *Room
	South *Room
	West  *Room
	East  *Room
	Up    *Room
	Down  *Room
}

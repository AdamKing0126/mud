package interfaces

import "database/sql"

type Area interface {
	GetUUID() string
	GetName() string
	GetDescription() string
	GetRooms() []Room
	GetRoomByUUID(string) (Room, error)
	SetRooms([]Room)
	SetRoomAtIndex(idx int, room Room)
	Run(db *sql.DB, ch chan Action, connections map[string]Player)
}

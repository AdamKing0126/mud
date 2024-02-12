package interfaces

import "github.com/jmoiron/sqlx"

type Area interface {
	GetUUID() string
	GetName() string
	GetDescription() string
	GetRooms() []Room
	GetRoomByUUID(string) (Room, error)
	SetRooms([]Room)
	SetRoomAtIndex(idx int, room Room)
	Run(db *sqlx.DB, ch chan Action, connections map[string]Player)
}

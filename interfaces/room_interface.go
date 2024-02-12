package interfaces

import "github.com/jmoiron/sqlx"

type Room interface {
	GetUUID() string
	GetName() string
	GetDescription() string
	GetPlayers() []Player
	GetPlayerByName(string) Player
	SetPlayers([]Player)
	GetMobs() []Mob
	SetMobs([]Mob)
	GetExits() ExitInfo
	SetItems([]Item)
	GetItems() []Item
	AddPlayer(Player)
	AddItem(*sqlx.DB, Item) error
	RemoveItem(Item) error
	RemovePlayer(Player) error
	SetExits(ExitInfo)
}

type ExitInfo interface {
	GetNorth() Room
	GetSouth() Room
	GetWest() Room
	GetEast() Room
	GetUp() Room
	GetDown() Room
}

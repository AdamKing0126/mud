package interfaces

import (
	"database/sql"
	"net"
)

type Player interface {
	AddItem(db *sql.DB, item Item) error
	RemoveItem(item Item) error
	DisplayEquipment()
	Equip(db *sql.DB, item Item) bool
	GetAbilities() Abilities
	GetArmorClass() int
	GetArea() Area
	GetAreaUUID() string
	GetColorProfile() ColorProfile
	GetColorProfileFromDB(db *sql.DB) error
	GetCommands() []string
	GetConn() net.Conn
	GetEquipment() PlayerEquipment
	GetEquipmentFromDB(db *sql.DB) error
	GetHashedPassword() string
	GetHealth() int
	GetHealthMax() int
	GetInventory() []Item
	GetItemFromInventory(string) Item
	GetLoggedIn() bool
	GetMana() int
	GetManaMax() int
	GetMovement() int
	GetMovementMax() int
	GetName() string
	GetRoom() Room
	GetRoomUUID() string
	GetUUID() string
	Logout(db *sql.DB) error
	Regen(db *sql.DB) error
	Remove(db *sql.DB, itemName string)
	SetAbilities(PlayerAbilities) error
	SetCommands([]string)
	SetConn(net.Conn)
	SetHealth(int)
	SetLocation(db *sql.DB, roomUUID string) error
	SetInventory([]Item)
	SetRoom(Room)
}

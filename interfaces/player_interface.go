package interfaces

import (
	"net"

	"github.com/jmoiron/sqlx"
)

type Player interface {
	AddItem(db *sqlx.DB, item Item) error
	RemoveItem(item Item) error
	DisplayEquipment()
	Equip(db *sqlx.DB, item Item) bool
	GetAbilities() Abilities
	GetArmorClass() int
	GetArea() Area
	GetAreaUUID() string
	GetColorProfile() ColorProfile
	GetColorProfileFromDB(db *sqlx.DB) error
	GetCommands() []string
	GetConn() net.Conn
	GetEquipment() PlayerEquipment
	GetEquipmentFromDB(db *sqlx.DB) error
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
	Logout(db *sqlx.DB) error
	Regen(db *sqlx.DB) error
	Remove(db *sqlx.DB, itemName string)
	SetAbilities(PlayerAbilities) error
	SetCommands([]string)
	SetConn(net.Conn)
	SetHealth(int)
	SetLocation(db *sqlx.DB, roomUUID string) error
	SetInventory([]Item)
	SetRoom(Room)
}

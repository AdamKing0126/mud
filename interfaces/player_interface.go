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
	GetArmorClass() int32
	GetArea() Area
	GetAreaUUID() string
	GetColorProfile() ColorProfile
	GetColorProfileFromDB(db *sqlx.DB) error
	GetCommands() []string
	GetConn() net.Conn
	GetEquipment() PlayerEquipment
	GetEquipmentFromDB(db *sqlx.DB) error
	GetHashedPassword() string
	GetHealth() int32
	GetHealthMax() int32
	GetInventory() []Item
	GetItemFromInventory(string) Item
	GetLoggedIn() bool
	GetMana() int32
	GetManaMax() int32
	GetMovement() int32
	GetMovementMax() int32
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
	SetHealth(int32)
	SetLocation(db *sqlx.DB, roomUUID string) error
	SetInventory([]Item)
	SetRoom(Room)
}

package interfaces

import (
	"database/sql"
	"net"
)

type Player interface {
	GetUUID() string
	GetName() string
	GetRoom() string
	GetArea() string
	GetHealth() int
	GetHealthMax() int
	GetMovement() int
	GetMovementMax() int
	GetMana() int
	GetManaMax() int
	SetHealth(int)
	GetConn() net.Conn
	SetConn(net.Conn)
	GetColorProfileFromDB(db *sql.DB) error
	GetEquipment() PlayerEquipment
	GetInventory() []Item
	GetEquipmentFromDB(db *sql.DB) error
	AddItemToInventory(db *sql.DB, item Item) error
	SetLocation(db *sql.DB, roomUUID string) error
	Logout(db *sql.DB) error
	GetCommands() []string
	SetCommands([]string)
	GetColorProfile() ColorProfile
	GetLoggedIn() bool
	Regen(db *sql.DB) error
	GetAbilities() Abilities
	GetArmorClass() int
	SetAbilities(PlayerAbilities)
	Equip(db *sql.DB, item Item) bool
	Remove(db *sql.DB, itemName string)
	GetHashedPassword() string
}

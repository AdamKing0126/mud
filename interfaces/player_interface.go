package interfaces

import (
	"database/sql"
	"net"
)

type PlayerInterface interface {
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
	GetEquipment() PlayerEquipmentInterface
	GetEquipmentFromDB(db *sql.DB) error
	SetLocation(db *sql.DB, roomUUID string) error
	Logout(db *sql.DB) error
	GetCommands() []string
	SetCommands([]string)
	GetColorProfile() ColorProfileInterface
	GetLoggedIn() bool
	Regen(db *sql.DB) error
	GetAbilities() AbilitiesInterface
	GetArmorClass() int
	SetAbilities(PlayerAbilitiesInterface)
	Equip(db *sql.DB, item Item) bool
	GetHashedPassword() string
}

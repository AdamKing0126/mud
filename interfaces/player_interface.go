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
	SetHealth(int)
	GetConn() net.Conn
	SetLocation(db *sql.DB, roomUUID string) error
	Logout(db *sql.DB) error
	GetCommands() []string
	SetCommands([]string)
	GetColorProfile() ColorProfileInterface
	GetLoggedIn() bool
}

type PlayerInRoomInterface interface {
	GetUUID() string
	GetName() string
}

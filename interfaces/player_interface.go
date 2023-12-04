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
	GetConn() net.Conn
	SetLocation(db *sql.DB, roomUUID string) error
	Logout(db *sql.DB) error
	GetCommands() []string
	SetCommands([]string)
	GetColorProfile() ColorProfileInterface
}

type PlayerInRoomInterface interface {
	GetUUID() string
	GetName() string
}

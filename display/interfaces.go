package display

import (
	"mud/interfaces"
	"net"
)

type Player interface {
	GetColorProfile() interfaces.ColorProfile
	GetConn() net.Conn
}

type ColorProfile interface {
	GetUUID() string
	GetColor(string) string
}

package display

import (
	"mud/interfaces"
	"net"
)

type Player interface {
	GetColorProfile() interfaces.ColorProfileInterface
	GetConn() net.Conn
}

type ColorProfile interface {
	GetUUID() string
	GetColor(string) string
}

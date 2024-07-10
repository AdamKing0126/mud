package display

import (
	"net"
)

type Player interface {
	GetColorProfile() ColorProfile
	GetColorProfileColor(string) string
	GetConn() net.Conn
}

type ColorProfile interface {
	GetUUID() string
	GetColor(string) string
}

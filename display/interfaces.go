package display

import "net"

type Player interface {
	GetColorProfile() ColorProfile
	GetConn() net.Conn
}

type ColorProfile interface {
	GetUUID() string
	GetColor(string) string
}

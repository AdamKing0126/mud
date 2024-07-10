package interfaces

import (
	"net"
)

type ProfilePlayer interface {
	GetColorProfileColor(string) string
	GetConn() net.Conn
}

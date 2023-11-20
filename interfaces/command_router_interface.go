package interfaces

import (
	"database/sql"
)

type CommandRouterInterface interface {
	HandleCommand(db *sql.DB, player PlayerInterface, command []byte)
}

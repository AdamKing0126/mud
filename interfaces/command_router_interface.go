package interfaces

import (
	"database/sql"
)

type CommandRouter interface {
	HandleCommand(db *sql.DB, player Player, command []byte)
}

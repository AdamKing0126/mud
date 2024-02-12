package interfaces

import "github.com/jmoiron/sqlx"

type CommandRouter interface {
	HandleCommand(db *sqlx.DB, player Player, command []byte)
}

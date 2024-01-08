package interfaces

import "database/sql"

type Area interface {
	GetUUID() string
	GetName() string
	GetDescription() string
	Run(db *sql.DB, ch chan Action, connections map[string]Player)
}

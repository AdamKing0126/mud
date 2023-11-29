package interfaces

import "database/sql"

type AreaInterface interface {
	GetUUID() string
	GetName() string
	GetDescription() string
	Run(db *sql.DB, ch chan ActionInterface)
}

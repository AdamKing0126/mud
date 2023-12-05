package interfaces

import "database/sql"

type ItemInterface interface {
	GetUUID() string
	GetName() string
	GetDescription() string
	SetLocation(db *sql.DB, playerUUID string, roomUUID string) error
}

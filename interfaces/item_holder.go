package interfaces

import "database/sql"

type ItemHolder interface {
	AddItem(db *sql.DB, item Item) error
	RemoveItem(item Item) error
}

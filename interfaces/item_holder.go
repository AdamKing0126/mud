package interfaces

import "github.com/jmoiron/sqlx"

type ItemHolder interface {
	AddItem(db *sqlx.DB, item Item) error
	RemoveItem(item Item) error
}

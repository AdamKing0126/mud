package players

import (
	"database/sql"
)

type Item interface {
	GetUUID() string
	GetName() string
	GetDescription() string
	SetLocation(db *sql.DB, playerUUID string, roomUUID string) error
	GetEquipmentSlots() []string
}

type EquippedItem interface {
	Item
	GetEquippedSlot() string
}

package interfaces

import "database/sql"

type EquipmentSlot string

type ItemInterface interface {
	GetUUID() string
	GetName() string
	GetDescription() string
	SetLocation(db *sql.DB, playerUUID string, roomUUID string) error
	GetEquipmentSlots() []EquipmentSlot
}

type EquippedItemInterface interface {
	ItemInterface
	GetEquippedSlot() string
}

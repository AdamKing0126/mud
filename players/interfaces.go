package players

import "github.com/jmoiron/sqlx"

type Item interface {
	GetUUID() string
	GetName() string
	GetDescription() string
	SetLocation(db *sqlx.DB, playerUUID string, roomUUID string) error
	GetEquipmentSlots() []string
}

type EquippedItem interface {
	Item
	GetEquippedSlot() string
}

type PlayerEquippedItem struct {
	Item
	EquippedSlot string
}

func (p *PlayerEquippedItem) GetEquippedSlot() string {
	return p.EquippedSlot
}

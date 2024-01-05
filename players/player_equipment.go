package players

import (
	"database/sql"
	"mud/interfaces"
)

func NewPlayerEquipment() *PlayerEquipment {
	return &PlayerEquipment{}
}

type PlayerEquipment struct {
	UUID         string
	PlayerUUID   string
	Head         interfaces.EquippedItemInterface
	Neck         interfaces.EquippedItemInterface
	Chest        interfaces.EquippedItemInterface
	Arms         interfaces.EquippedItemInterface
	Hands        interfaces.EquippedItemInterface
	DominantHand interfaces.EquippedItemInterface
	OffHand      interfaces.EquippedItemInterface
	Legs         interfaces.EquippedItemInterface
	Feet         interfaces.EquippedItemInterface
}

// TODO working on this, I can't remember at the moment why i left this in here
func (pe *PlayerEquipment) GetEquippedLocation(db *sql.DB, item interfaces.EquippedItemInterface) string {
	return "foo"
}

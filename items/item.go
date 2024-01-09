package items

import (
	"github.com/google/uuid"
)

const (
	Head         = "Head"
	Neck         = "Neck"
	Chest        = "Chest"
	Arms         = "Arms"
	Hands        = "Hands"
	DominantHand = "DominantHand"
	OffHand      = "OffHand"
	Legs         = "Legs"
	Feet         = "Feet"
)

func NewItem(uuid, name, description string, equipmentSlots []string) *Item {
	return &Item{
		UUID:           uuid,
		Name:           name,
		Description:    description,
		EquipmentSlots: equipmentSlots,
	}
}

func NewEquippedItem(item *Item, equippedSlot string) *EquippedItem {
	return &EquippedItem{
		Item:         item,
		EquippedSlot: equippedSlot,
	}
}

type Item struct {
	UUID           string
	Name           string
	Description    string
	EquipmentSlots []string
}

func (item *Item) GetUUID() string {
	return item.UUID
}

func (item *Item) GetName() string {
	return item.Name
}

func (item *Item) GetDescription() string {
	return item.Description
}

func (item *Item) GetEquipmentSlots() []string {
	return item.EquipmentSlots
}

type EquippedItem struct {
	*Item
	EquippedSlot string
}

func (ei EquippedItem) GetEquippedSlot() string {
	return ei.EquippedSlot
}

type ItemLocation struct {
	ItemUUID   uuid.UUID
	RoomUUID   uuid.UUID
	PlayerUUID uuid.UUID
}

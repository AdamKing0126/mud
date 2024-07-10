package players

import (
	"mud/items"

	"github.com/jmoiron/sqlx"
)

func NewPlayerEquipment() *PlayerEquipment {
	return &PlayerEquipment{}
}

type PlayerEquipment struct {
	UUID         string
	PlayerUUID   string
	Head         *EquippedItem
	Neck         *EquippedItem
	Chest        *EquippedItem
	Arms         *EquippedItem
	Hands        *EquippedItem
	DominantHand *EquippedItem
	OffHand      *EquippedItem
	Legs         *EquippedItem
	Feet         *EquippedItem
}

type EquippedItem struct {
	items.Item
	EquippedSlot string
}

func NewEquippedItem(item items.Item, equippedSlot string) *EquippedItem {
	return &EquippedItem{
		Item:         item,
		EquippedSlot: equippedSlot,
	}
}

func (pe PlayerEquipment) GetUUID() string {
	return pe.UUID
}

func (pe PlayerEquipment) GetPlayerUUID() string {
	return pe.PlayerUUID
}

func (pe PlayerEquipment) GetHead() *EquippedItem {
	return pe.Head
}

func (pe PlayerEquipment) GetNeck() *EquippedItem {
	return pe.Neck
}

func (pe PlayerEquipment) GetChest() *EquippedItem {
	return pe.Chest
}

func (pe PlayerEquipment) GetArms() *EquippedItem {
	return pe.Arms
}

func (pe PlayerEquipment) GetHands() *EquippedItem {
	return pe.Hands
}

func (pe PlayerEquipment) GetDominantHand() *EquippedItem {
	return pe.DominantHand
}

func (pe PlayerEquipment) GetOffHand() *EquippedItem {
	return pe.OffHand
}

func (pe PlayerEquipment) GetLegs() *EquippedItem {
	return pe.Legs
}

func (pe PlayerEquipment) GetFeet() *EquippedItem {
	return pe.Feet
}

// TODO working on this, I can't remember at the moment why i left this in here
func (pe *PlayerEquipment) GetEquippedLocation(db *sqlx.DB, item EquippedItem) string {
	return "foo"
}

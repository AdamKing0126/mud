package players

import (
	"database/sql"
	"fmt"
	"mud/interfaces"
)

func NewPlayerEquipment() *PlayerEquipment {
	return &PlayerEquipment{}
}

type PlayerEquipment struct {
	UUID         string
	PlayerUUID   string
	Head         EquippedItem
	Neck         EquippedItem
	Chest        EquippedItem
	Arms         EquippedItem
	Hands        EquippedItem
	DominantHand EquippedItem
	OffHand      EquippedItem
	Legs         EquippedItem
	Feet         EquippedItem
}

func (pe *PlayerEquipment) GetUUID() string {
	return pe.UUID
}

func (pe *PlayerEquipment) GetPlayerUUID() string {
	return pe.PlayerUUID
}

func (pe *PlayerEquipment) GetHead() interfaces.EquippedItem {
	return pe.Head
}

func (pe *PlayerEquipment) GetNeck() interfaces.EquippedItem {
	return pe.Neck
}

func (pe *PlayerEquipment) GetChest() interfaces.EquippedItem {
	return pe.Chest
}

func (pe *PlayerEquipment) GetArms() interfaces.EquippedItem {
	return pe.Arms
}

func (pe *PlayerEquipment) GetHands() interfaces.EquippedItem {
	return pe.Hands
}

func (pe *PlayerEquipment) GetDominantHand() interfaces.EquippedItem {
	return pe.DominantHand
}

func (pe *PlayerEquipment) SetDominantHand(item *Item) {
	fmt.Println("poo")
}

func (pe *PlayerEquipment) GetOffHand() interfaces.EquippedItem {
	return pe.OffHand
}

func (pe *PlayerEquipment) GetLegs() interfaces.EquippedItem {
	return pe.Legs
}

func (pe *PlayerEquipment) GetFeet() interfaces.EquippedItem {
	return pe.Feet
}

// TODO working on this, I can't remember at the moment why i left this in here
func (pe *PlayerEquipment) GetEquippedLocation(db *sql.DB, item EquippedItem) string {
	return "foo"
}

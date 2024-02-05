package items

import (
	"database/sql"
	"encoding/json"
	"fmt"

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

func NewItemFromTemplate(db *sql.DB, templateUUID string) (*Item, error) {
	var name, description string
	var equipmentSlotsJSON string
	query := `SELECT name, description, equipment_slots
				FROM item_templates 
				WHERE uuid = ?`
	err := db.QueryRow(query, templateUUID).Scan(&name, &description, &equipmentSlotsJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %v", err)
	}

	itemUUID := uuid.NewString()
	query = `INSERT INTO items (uuid, name, description, equipment_slots)
				VALUES (?, ?, ?, ?)`
	_, err = db.Exec(query, itemUUID, name, description, equipmentSlotsJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to insert item: %v", err)
	}

	var equipmentSlots []string
	err = json.Unmarshal([]byte(equipmentSlotsJSON), &equipmentSlots)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %v", err)
	}

	fmt.Printf("creating item with equipmentSlots %v", equipmentSlots)
	return NewItem(itemUUID, name, description, equipmentSlots), nil
}

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

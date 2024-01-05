package items

import (
	"database/sql"
	"fmt"
	"strings"

	"mud/interfaces"

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

type ItemLocation struct {
	ItemUUID   uuid.UUID
	RoomUUID   uuid.UUID
	PlayerUUID uuid.UUID
}

type EquippedItem struct {
	Item
	EquippedSlot string
}

func NewEquippedItem(item interfaces.ItemInterface, equippedSlot string) interfaces.EquippedItemInterface {
	itemConcrete, ok := item.(*Item)
	if !ok {
		return nil
	}
	return &EquippedItem{
		Item:         *itemConcrete,
		EquippedSlot: equippedSlot,
	}
}

func (ei EquippedItem) GetEquippedSlot() string {
	return ei.EquippedSlot
}

type Item struct {
	UUID           string
	Name           string
	Description    string
	EquipmentSlots []interfaces.EquipmentSlot
}

func NewItem(uuid, name, description string, equipmentSlots []string) interfaces.ItemInterface {
	e := make([]interfaces.EquipmentSlot, len(equipmentSlots))
	for i, v := range e {
		e[i] = interfaces.EquipmentSlot(v)
	}
	return &Item{
		UUID:           uuid,
		Name:           name,
		Description:    description,
		EquipmentSlots: e,
	}
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

func (item *Item) GetEquipmentSlots() []interfaces.EquipmentSlot {
	return item.EquipmentSlots
}

func (item *Item) SetLocation(db *sql.DB, playerUUID string, roomUUID string) error {
	var query string
	if playerUUID != "" {
		query = fmt.Sprintf("UPDATE item_locations SET room_uuid = '', player_uuid = '%s' WHERE item_uuid = '%s'", playerUUID, item.UUID)
	} else {
		query = fmt.Sprintf("UPDATE item_locations SET room_uuid = '%s', player_uuid = '' WHERE item_uuid = '%s'", roomUUID, item.UUID)
	}
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func GetEquippedItemsForPlayer(db *sql.DB, playerUUID string) ([]interfaces.EquippedItemInterface, error) {
	query := `
	SELECT i.uuid, i.name,
	CASE
		WHEN pe.Head = i.uuid THEN 'Head'
		WHEN pe.Neck = i.uuid THEN 'Neck'
		WHEN pe.Chest = i.uuid THEN 'Chest'
		WHEN pe.Arms = i.uuid THEN 'Arms'
		WHEN pe.Hands = i.uuid THEN 'Hands'
		WHEN pe.DominantHand = i.uuid THEN 'DominantHand'
		WHEN pe.Legs = i.uuid THEN 'Legs'
		WHEN pe.Feet = i.uuid THEN 'Feet'
		ELSE NULL
	END AS equipped_slot
	FROM item_locations il
	JOIN items i on il.item_uuid = i.uuid
	JOIN player_equipments pe on pe.player_uuid = il.player_uuid
	`

	rows, err := db.Query(query, playerUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	var items []EquippedItem
	for rows.Next() {
		var item EquippedItem
		err := rows.Scan(&item.UUID, &item.Name, &item.EquippedSlot)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		items = append(items, item)
	}

	itemInterfaces := make([]interfaces.EquippedItemInterface, len(items))
	for i := range items {
		itemInterfaces[i] = &items[i]
	}
	return itemInterfaces, nil
}

func GetItemsForPlayer(db *sql.DB, playerUUID string) ([]interfaces.ItemInterface, error) {
	query := `
		SELECT i.uuid, i.name, i.description, i.equipment_slots
		FROM item_locations il
		JOIN items i ON il.item_uuid = i.uuid
		JOIN player_equipments pe ON pe.player_uuid = il.player_uuid
		WHERE il.player_uuid = ? AND NOT EXISTS (
			SELECT 1
			FROM player_equipments pe2
			WHERE pe2.player_uuid = il.player_uuid
			AND (pe2.Head = i.uuid OR pe2.Neck = i.uuid OR pe2.Chest = i.uuid OR pe2.Arms = i.uuid OR pe2.Hands = i.uuid OR pe2.DominantHand = i.uuid OR pe2.OffHand = i.uuid OR pe2.Legs = i.uuid OR pe2.Feet = i.uuid)
		)
	`

	rows, err := db.Query(query, playerUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	var items []Item
	for rows.Next() {
		var item Item
		var equipmentSlots string
		err := rows.Scan(&item.UUID, &item.Name, &item.Description, &equipmentSlots)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		for _, slot := range strings.Split(equipmentSlots, ",") {
			switch slot {
			case "Head":
				item.EquipmentSlots = append(item.EquipmentSlots, Head)
			case "Neck":
				item.EquipmentSlots = append(item.EquipmentSlots, Neck)
			case "Chest":
				item.EquipmentSlots = append(item.EquipmentSlots, Chest)
			case "Arms":
				item.EquipmentSlots = append(item.EquipmentSlots, Arms)
			case "Hands":
				item.EquipmentSlots = append(item.EquipmentSlots, Hands)
			case "DominantHand":
				item.EquipmentSlots = append(item.EquipmentSlots, DominantHand)
			case "OffHand":
				item.EquipmentSlots = append(item.EquipmentSlots, OffHand)
			case "Legs":
				item.EquipmentSlots = append(item.EquipmentSlots, Legs)
			case "Feet":
				item.EquipmentSlots = append(item.EquipmentSlots, Feet)
			}
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	itemInterfaces := make([]interfaces.ItemInterface, len(items))
	for i, item := range items {
		itemInterfaces[i] = &item
	}
	return itemInterfaces, nil
}

func GetItemByNameForPlayer(db *sql.DB, itemName string, playerUUID string) (interfaces.ItemInterface, error) {
	items, err := GetItemsForPlayer(db, playerUUID)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.GetName() == itemName {
			return item, nil
		}
	}

	return nil, fmt.Errorf("item not found")
}

func GetItemsInRoom(db *sql.DB, roomUUID string) ([]Item, error) {
	query := `
		SELECT i.uuid, i.name, i.description
		FROM item_locations il
		JOIN items i ON il.item_uuid = i.uuid
		WHERE il.room_uuid = ?
	`
	rows, err := db.Query(query, roomUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.UUID, &item.Name, &item.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return items, nil
}

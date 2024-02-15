package items

import (
	"encoding/json"
	"fmt"
	"mud/interfaces"
	"strings"

	"github.com/jmoiron/sqlx"
)

func GetItemsInRoom(db *sqlx.DB, roomUUID string) ([]interfaces.Item, error) {
	query := `
		SELECT i.uuid, i.name, i.description, i.equipment_slots
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
		var equipmentSlotsJSON string
		var equipmentSlots []string
		err := rows.Scan(&item.UUID, &item.Name, &item.Description, &equipmentSlotsJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		err = json.Unmarshal([]byte(equipmentSlotsJSON), &equipmentSlots)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
		}

		for _, slot := range equipmentSlots {
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
	itemInterfaces := make([]interfaces.Item, len(items))
	for idx := range items {
		itemInterfaces[idx] = &items[idx]
	}

	return itemInterfaces, nil
}

func (item *Item) SetLocation(db *sqlx.DB, playerUUID string, roomUUID string) error {
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

func GetItemsForPlayer(db *sqlx.DB, playerUUID string) ([]interfaces.Item, error) {
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

	var itemInterfaces []interfaces.Item
	for _, item := range items {
		itemInterfaces = append(itemInterfaces, &item)
	}

	return itemInterfaces, nil
}

func GetItemByNameForPlayer(db *sqlx.DB, itemName string, playerUUID string) (interfaces.Item, error) {
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

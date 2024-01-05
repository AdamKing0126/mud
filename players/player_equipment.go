package players

import (
	"database/sql"
	"fmt"
	"mud/interfaces"
	"mud/items"
	"strings"
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

func GetPlayerEquipment(db *sql.DB, playerUUID string) (*PlayerEquipment, error) {
	var head, neck, chest, arms, hands, dominantHand, offHand, legs, feet string
	query := `SELECT uuid, player_uuid, Head, Neck, Chest, Arms, Hands, DominantHand, OffHand, Legs, Feet 
			  FROM player_equipments 
			  WHERE player_uuid = ?`

	var pe PlayerEquipment
	err := db.QueryRow(query, playerUUID).Scan(&pe.UUID, &pe.PlayerUUID, &head, &neck, &chest, &arms, &hands, &dominantHand, &offHand, &legs, &feet)

	item_uuids := []string{head, neck, chest, arms, hands, dominantHand, offHand, legs, feet}
	// itemsMap := make(map[string]interfaces.EquippedItemInterface)

	placeholders := strings.Trim(strings.Repeat("?,", len(item_uuids)), ",")

	queryString := fmt.Sprintf("SELECT uuid, name, description, equipment_slots FROM items where uuid in (%s)", placeholders)

	args := make([]interface{}, len(item_uuids))
	for i, v := range item_uuids {
		args[i] = v
	}

	rows, err := db.Query(queryString, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var uuid, name, description, equipmentSlots string
		err := rows.Scan(&uuid, &name, &description, &equipmentSlots)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		equipmentSlotArray := strings.Split(equipmentSlots, ",")
		item := items.NewItem(uuid, name, description, equipmentSlotArray)

		switch uuid {
		case head:
			pe.Head = items.NewEquippedItem(item, "Head")
		case neck:
			pe.Neck = items.NewEquippedItem(item, "Neck")
		case chest:
			pe.Chest = items.NewEquippedItem(item, "Chest")
		case arms:
			pe.Arms = items.NewEquippedItem(item, "Arms	")
		case hands:
			pe.Hands = items.NewEquippedItem(item, "Hands")
		case dominantHand:
			pe.DominantHand = items.NewEquippedItem(item, "DominantHand")
		case offHand:
			pe.OffHand = items.NewEquippedItem(item, "OffHand")
		case legs:
			pe.Legs = items.NewEquippedItem(item, "Legs")
		default:
			pe.Feet = items.NewEquippedItem(item, "Feet")
		}
	}

	if err != nil {
		return nil, err
	}

	return &pe, nil
}

func (pe *PlayerEquipment) GetEquippedLocation(db *sql.DB, item interfaces.EquippedItemInterface) string {
	return "foo"
}

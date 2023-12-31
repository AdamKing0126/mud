package players

import (
	"database/sql"
	"fmt"
	"mud/interfaces"
	"mud/items"
	"strings"
)

func GetPlayerFromDB(db *sql.DB, playerName string) (interfaces.PlayerInterface, error) {
	var player Player
	var colorProfileUUID string
	err := db.QueryRow("SELECT uuid, name, room, area, health, health_max, movement, movement_max, mana, mana_max, logged_in, password, color_profile FROM players WHERE LOWER(name) = LOWER(?)", playerName).
		Scan(&player.UUID, &player.Name, &player.Room, &player.Area, &player.Health, &player.HealthMax, &player.Movement, &player.MovementMax, &player.Mana, &player.ManaMax, &player.LoggedIn, &player.Password, &colorProfileUUID)
	if err != nil {
		return nil, err
	}
	player.ColorProfile = &ColorProfile{UUID: colorProfileUUID}

	return &player, nil
}

func (player *Player) GetColorProfileFromDB(db *sql.DB) error {
	var colorProfile = &ColorProfile{}
	query := `SELECT uuid, name, primary_color, secondary_color, warning_color, danger_color, title_color, description_color 
	FROM color_profiles WHERE uuid = ?;`
	err := db.QueryRow(query, player.GetColorProfile().GetUUID()).Scan(&colorProfile.UUID, &colorProfile.Name, &colorProfile.Primary, &colorProfile.Secondary, &colorProfile.Warning, &colorProfile.Danger, &colorProfile.Title, &colorProfile.Description)
	if err != nil {
		return err
	}
	player.ColorProfile = colorProfile
	return nil
}

func setPlayerLoggedInStatusInDB(db *sql.DB, playerUUID string, loggedIn bool) error {
	_, err := db.Exec("UPDATE players SET logged_in = ? WHERE uuid = ?", loggedIn, playerUUID)
	if err != nil {
		return err
	}
	return nil
}

func (player *Player) GetEquipmentFromDB(db *sql.DB) error {
	var head, neck, chest, arms, hands, dominantHand, offHand, legs, feet string
	query := `SELECT uuid, player_uuid, Head, Neck, Chest, Arms, Hands, DominantHand, OffHand, Legs, Feet 
			  FROM player_equipments 
			  WHERE player_uuid = ?`

	var pe PlayerEquipment
	err := db.QueryRow(query, player.UUID).Scan(&pe.UUID, &pe.PlayerUUID, &head, &neck, &chest, &arms, &hands, &dominantHand, &offHand, &legs, &feet)
	if err != nil {
		return err
	}

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
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var uuid, name, description, equipmentSlots string
		err := rows.Scan(&uuid, &name, &description, &equipmentSlots)
		if err != nil {
			return fmt.Errorf("failed to scan row: %v", err)
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
		return err
	}

	player.Equipment = pe

	return nil
}

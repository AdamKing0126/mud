package players

import (
	"context"
	"fmt"
	"strings"

	"github.com/adamking0126/mud/internal/game/character_classes"
	"github.com/adamking0126/mud/internal/game/items"
	"github.com/adamking0126/mud/pkg/database"
)

type Repository struct {
	db database.DB
}

func NewRepository(db database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SetPlayerLoggedInStatus(ctx context.Context, playerUUID string, loggedIn bool) error {
	err := r.db.Exec(ctx, "UPDATE players SET logged_in = ? WHERE uuid = ?", loggedIn, playerUUID)
	return err
}

func (r *Repository) GetPlayerByName(ctx context.Context, name string) (*Player, error) {
	return GetPlayerByName(ctx, r.db, name)
}

func GetPlayerByName(ctx context.Context, db database.DB, name string) (*Player, error) {
	var player Player
	var playerAbilities PlayerAbilities
	db.QueryRow(ctx, `
		SELECT p.uuid, p.name, p.room, p.area, p.hp, p.movement, p.logged_in, 
		       pa.intelligence, pa.dexterity, pa.charisma, pa.constitution, pa.wisdom, pa.strength 
		FROM players p 
		JOIN player_attributes pa ON p.uuid = pa.player_uuid 
		WHERE LOWER(p.name) = LOWER(?)`, name).Scan(&player.UUID, &player.Name, &player.RoomUUID, &player.AreaUUID, &player.HP, &player.Movement, &player.LoggedIn,
		&playerAbilities.Intelligence, &playerAbilities.Dexterity, &playerAbilities.Charisma,
		&playerAbilities.Constitution, &playerAbilities.Wisdom, &playerAbilities.Strength)

	// TODO ADAM
	// player.Abilities = playerAbilities
	return &player, nil
}

func (r *Repository) GetPlayerByNameFull(ctx context.Context, playerName string) (*Player, error) {
	var player Player
	var colorProfileUUID string
	var characterClassArchetypeSlug string
	var characterRaceSlug, characterSubRaceSlug string
	r.db.QueryRow(ctx, `
		SELECT uuid, name, character_class, race, subrace, room, area, hp, hp_max, movement, movement_max, logged_in, password, color_profile 
		FROM players 
		WHERE LOWER(name) = LOWER(?)`, playerName).Scan(&player.UUID, &player.Name, &characterClassArchetypeSlug, &characterRaceSlug, &characterSubRaceSlug,
		&player.RoomUUID, &player.AreaUUID, &player.HP, &player.HPMax, &player.Movement, &player.MovementMax,
		&player.LoggedIn, &player.Password, &colorProfileUUID)

	// Note: might want to move these to separate repository methods
	characterClasses, err := character_classes.GetCharacterClassList(ctx, r.db, characterClassArchetypeSlug)
	if err != nil || len(characterClasses) != 1 {
		return nil, err
	}
	player.CharacterClass = characterClasses[0]

	characterRaces, err := character_classes.GetCharacterRaceList(ctx, r.db, characterRaceSlug, characterSubRaceSlug)
	if err != nil || len(characterRaces) != 1 {
		return nil, err
	}
	player.Race = characterRaces[0]

	player.ColorProfile = ColorProfile{UUID: colorProfileUUID}

	return &player, nil
}

func (r *Repository) GetColorProfile(ctx context.Context, uuid string) *ColorProfile {
	var colorProfile ColorProfile
	query := `SELECT uuid, name, primary_color, secondary_color, warning_color, danger_color, title_color, description_color 
			  FROM color_profiles WHERE uuid = ?`
	r.db.QueryRow(ctx, query, uuid).Scan(
		&colorProfile.UUID, &colorProfile.Name, &colorProfile.Primary, &colorProfile.Secondary,
		&colorProfile.Warning, &colorProfile.Danger, &colorProfile.Title, &colorProfile.Description)
	return &colorProfile
}

func (r *Repository) GetPlayerInventory(ctx context.Context, playerUUID string) ([]*items.Item, error) {
	queryString := `SELECT i.uuid, i.name, i.description, i.equipment_slots 
					FROM item_locations il 
					JOIN items i ON il.player_uuid = ? AND il.item_uuid = i.uuid`
	rows, err := r.db.Query(ctx, queryString, playerUUID)
	if err != nil {
		return nil, fmt.Errorf("error querying inventory: %v", err)
	}
	defer rows.Close()

	var inventory []*items.Item
	for rows.Next() {
		var uuid, name, description, slots string
		err := rows.Scan(&uuid, &name, &description, &slots)
		if err != nil {
			return nil, fmt.Errorf("error scanning inventory row: %v", err)
		}
		equipmentSlots := strings.Split(slots, ",")
		item := items.NewItem(uuid, name, description, equipmentSlots)
		inventory = append(inventory, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating inventory rows: %v", err)
	}

	return inventory, nil
}

// func (r *Repository) GetPlayerEquipment(ctx context.Context, playerUUID string) (*PlayerEquipment, error) {
// ... (implement GetPlayerEquipment similar to GetEquipmentFromDB)
// }

func GetPlayersInRoom(ctx context.Context, db database.DB, roomUUID string) ([]*Player, error) {
	var players []*Player
	query := `SELECT uuid, name FROM players WHERE room = ? and logged_in = 1`
	rows, err := db.Query(ctx, query, roomUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		player := &Player{}
		err := rows.Scan(&player.UUID, &player.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		players = append(players, player)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return players, nil
}

func (player *Player) GetColorProfileFromDB(ctx context.Context, db database.DB) error {
	var colorProfile = ColorProfile{}
	query := `SELECT uuid, name, primary_color, secondary_color, warning_color, danger_color, title_color, description_color FROM color_profiles WHERE uuid = ?`
	err := db.QueryRow(ctx, query, player.ColorProfile.UUID).Scan(&colorProfile.UUID, &colorProfile.Name, &colorProfile.Primary, &colorProfile.Secondary, &colorProfile.Warning, &colorProfile.Danger, &colorProfile.Title, &colorProfile.Description)
	if err != nil {
		return err
	}
	player.ColorProfile = colorProfile
	return nil
}

func GetPlayerFromDB(ctx context.Context, db database.DB, playerName string) (*Player, error) {
	var player Player
	var colorProfileUUID string
	var characterClassArchetypeSlug string
	var characterRaceSlug, characterSubRaceSlug string
	err := db.QueryRow(ctx, "SELECT uuid, name, character_class, race, subrace, room, area, hp, hp_max, movement, movement_max, logged_in, password, color_profile FROM players WHERE LOWER(name) = LOWER(?)", playerName).
		Scan(&player.UUID, &player.Name, &characterClassArchetypeSlug, &characterRaceSlug, &characterSubRaceSlug, &player.RoomUUID, &player.AreaUUID, &player.HP, &player.HPMax, &player.Movement, &player.MovementMax, &player.LoggedIn, &player.Password, &colorProfileUUID)
	if err != nil {
		return nil, err
	}
	return &player, nil
}

// package players

// import (
// 	"fmt"

// 	"github.com/adamking0126/mud/internal/game/character_classes"
// 	"github.com/adamking0126/mud/internal/game/items"

// 	"strings"

// 	"github.com/jmoiron/sqlx"
// )

func setPlayerLoggedInStatusInDB(ctx context.Context, db database.DB, playerUUID string, loggedIn bool) error {
	err := db.Exec(ctx, "UPDATE players SET logged_in = ? WHERE uuid = ?", loggedIn, playerUUID)
	if err != nil {
		return err
	}
	return nil
}

// func GetPlayerByName(db *sqlx.DB, name string) (*Player, error) {
// 	var player Player
// 	var playerAbilities PlayerAbilities
// 	err := db.QueryRow("SELECT p.uuid, p.name, p.room, p.area, p.hp, p.movement, p.logged_in, pa.intelligence, pa.dexterity, pa.charisma, pa.constitution, pa.wisdom, pa.strength FROM players p JOIN player_attributes pa ON p.uuid = pa.player_uuid WHERE LOWER(p.name) = LOWER(?)", name).
// 		Scan(&player.UUID, &player.Name, &player.RoomUUID, &player.AreaUUID, &player.HP, &player.Movement, &player.LoggedIn, &playerAbilities.Intelligence, &playerAbilities.Dexterity, &playerAbilities.Charisma, &playerAbilities.Constitution, &playerAbilities.Wisdom, &playerAbilities.Strength)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &player, nil
// }

// 	characterClasses, err := character_classes.GetCharacterClassList(db, characterClassArchetypeSlug)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(characterClasses) != 1 {
// 		return nil, nil
// 	}
// 	var characterClass = characterClasses[0]
// 	player.CharacterClass = characterClass

// 	characterRaces, err := character_classes.GetCharacterRaceList(db, characterRaceSlug, characterSubRaceSlug)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(characterRaces) != 1 {
// 		return nil, nil
// 	}
// 	player.Race = characterRaces[0]

// 	player.ColorProfile = ColorProfile{UUID: colorProfileUUID}

// 	return &player, nil
// }

func (player *Player) GetInventoryFromDB(ctx context.Context, db database.DB) error {
	queryString := `SELECT i.uuid, i.name, i.description, i.equipment_slots FROM item_locations il JOIN items i ON il.player_uuid = ? AND il.item_uuid = i.uuid;`
	rows, err := db.Query(ctx, queryString, player.UUID)
	if err != nil {
		fmt.Printf("error querying: %v", err)
	}
	var inventory []*items.Item
	for rows.Next() {
		var uuid, name, description, slots string
		err := rows.Scan(&uuid, &name, &description, &slots)
		equipmentSlots := strings.Split(slots, ",")
		if err == nil {
			fmt.Println("uh oh")
		}
		item := items.NewItem(uuid, name, description, equipmentSlots)
		inventory = append(inventory, item)
	}

	player.Inventory = inventory
	return nil

}

func (player *Player) GetEquipmentFromDB(ctx context.Context, db database.DB) error {
	var head, neck, chest, arms, hands, dominantHand, offHand, legs, feet string
	query := `SELECT uuid, player_uuid, Head, Neck, Chest, Arms, Hands, DominantHand, OffHand, Legs, Feet
			  FROM player_equipments
			  WHERE player_uuid = ?`

	var pe PlayerEquipment
	err := db.QueryRow(ctx, query, player.UUID).Scan(&pe.UUID, &pe.PlayerUUID, &head, &neck, &chest, &arms, &hands, &dominantHand, &offHand, &legs, &feet)
	if err != nil {
		return err
	}
	player.Equipment = pe
	return nil
}

// 	item_uuids := []string{head, neck, chest, arms, hands, dominantHand, offHand, legs, feet}

// 	placeholders := strings.Trim(strings.Repeat("?,", len(item_uuids)), ",")

// 	queryString := fmt.Sprintf("SELECT uuid, name, description, equipment_slots FROM items where uuid in (%s)", placeholders)

// 	args := make([]interface{}, len(item_uuids))
// 	for i, v := range item_uuids {
// 		args[i] = v
// 	}

// 	rows, err := db.Query(queryString, args...)
// 	if err != nil {
// 		return err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var uuid, name, description, equipmentSlots string
// 		err := rows.Scan(&uuid, &name, &description, &equipmentSlots)
// 		if err != nil {
// 			return fmt.Errorf("failed to scan row: %v", err)
// 		}
// 		equipmentSlotArray := strings.Split(equipmentSlots, ",")
// 		item := items.NewItem(uuid, name, description, equipmentSlotArray)

// 		switch uuid {
// 		case head:
// 			pe.Head = NewEquippedItem(item, "Head")
// 		case neck:
// 			pe.Neck = NewEquippedItem(item, "Neck")
// 		case chest:
// 			pe.Chest = NewEquippedItem(item, "Chest")
// 		case arms:
// 			pe.Arms = NewEquippedItem(item, "Arms")
// 		case hands:
// 			pe.Hands = NewEquippedItem(item, "Hands")
// 		case dominantHand:
// 			pe.DominantHand = NewEquippedItem(item, "DominantHand")
// 		case offHand:
// 			pe.OffHand = NewEquippedItem(item, "OffHand")
// 		case legs:
// 			pe.Legs = NewEquippedItem(item, "Legs")
// 		default:
// 			pe.Feet = NewEquippedItem(item, "Feet")
// 		}
// 	}

// 	player.Equipment = pe

// 	return nil
// }

// 	return &colorProfile, nil
// }

// func GetPlayersInRoom(db *sqlx.DB, roomUUID string) ([]*Player, error) {
// 	// Would it be better to rely on the `connections` structure attached to the server
// 	// or is it better to query the db for this info?
// 	var players []*Player
// 	query := `
// 		SELECT uuid, name
// 		FROM players
// 		WHERE room = ? and logged_in = 1
// 	`
// 	rows, err := db.Query(query, roomUUID)
// 	if err != nil {
// 		return players, fmt.Errorf("failed to execute query: %v", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		player := &Player{}
// 		err := rows.Scan(&player.UUID, &player.Name)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to scan row: %v", err)
// 		}
// 		players = append(players, player)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("error iterating over rows: %v", err)
// 	}

// 	return players, nil
// }

func getColorProfileFromDB(ctx context.Context, db database.DB, uuid string) (*ColorProfile, error) {
	var colorProfile ColorProfile
	query := `SELECT uuid, name, primary_color, secondary_color, warning_color, danger_color, title_color, description_color FROM color_profiles WHERE uuid = ?`
	db.QueryRow(ctx, query, uuid).Scan(&colorProfile.UUID, &colorProfile.Name, &colorProfile.Primary, &colorProfile.Secondary, &colorProfile.Warning, &colorProfile.Danger, &colorProfile.Title, &colorProfile.Description)
	return &colorProfile, nil
}

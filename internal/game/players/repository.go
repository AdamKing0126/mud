package players

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/adamking0126/mud/internal/game/character_classes"
	"github.com/adamking0126/mud/internal/game/items"
	"github.com/adamking0126/mud/pkg/database"
	"github.com/charmbracelet/ssh"
	"github.com/google/uuid"
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
	var player Player
	var playerAbilities PlayerAbilities
	r.db.QueryRow(ctx, `
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

func (r *Repository) GetColorProfileForPlayerByUUID(ctx context.Context, uuid string) *ColorProfile {
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

func (r *Repository) GetPlayerFromDB(ctx context.Context, playerName string) (*Player, error) {
	var player Player
	var colorProfileUUID string
	var characterClassArchetypeSlug string
	var characterRaceSlug, characterSubRaceSlug string
	err := r.db.QueryRow(ctx, "SELECT uuid, name, character_class, race, subrace, room, area, hp, hp_max, movement, movement_max, logged_in, password, color_profile FROM players WHERE LOWER(name) = LOWER(?)", playerName).
		Scan(&player.UUID, &player.Name, &characterClassArchetypeSlug, &characterRaceSlug, &characterSubRaceSlug, &player.RoomUUID, &player.AreaUUID, &player.HP, &player.HPMax, &player.Movement, &player.MovementMax, &player.LoggedIn, &player.Password, &colorProfileUUID)
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *Repository) GetEquipmentForPlayerByUUID(ctx context.Context, playerUUID string) *PlayerEquipment {
	var head, neck, chest, arms, hands, dominantHand, offHand, legs, feet string
	query := `SELECT uuid, player_uuid, Head, Neck, Chest, Arms, Hands, DominantHand, OffHand, Legs, Feet
			  FROM player_equipments
			  WHERE player_uuid = ?`

	var pe PlayerEquipment
	err := r.db.QueryRow(ctx, query, playerUUID).Scan(&pe.UUID, &pe.PlayerUUID, &head, &neck, &chest, &arms, &hands, &dominantHand, &offHand, &legs, &feet)
	if err != nil {
		return nil
	}
	return &pe
}

func (r *Repository) GetInventoryForPlayerByUUID(ctx context.Context, playerUUID string) []*items.Item {
	queryString := `SELECT i.uuid, i.name, i.description, i.equipment_slots FROM item_locations il JOIN items i ON il.player_uuid = ? AND il.item_uuid = i.uuid;`
	rows, err := r.db.Query(ctx, queryString, playerUUID)
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
	return inventory
}

func (r *Repository) CreatePlayer(ctx context.Context, session ssh.Session, playerName string) (*Player, error) {
	player := NewPlayer(session)
	player.Name = playerName

	fmt.Fprintf(session, "Please enter a password you'd like to use: ")
	password := getPlayerInput(session)
	player.Password = HashPassword(password)

	// default start point
	player.AreaUUID = "d71e8cf1-d5ba-426c-8915-4c7f5b22e3a9"
	player.RoomUUID = "189a729d-4e40-4184-a732-e2c45c66ff46"
	player.UUID = uuid.New().String()

	// "default" light mode color profile.  Should let the user choose?
	colorProfile, err := r.GetColorProfileFromDB(ctx, "2c7dfd5b-d160-42e0-accb-b77d9686dbea")
	if err != nil {
		return nil, err
	}

	chosenCharacterClass := selectCharacterClassAndArchetype(ctx, session, r.db, player)
	if chosenCharacterClass == nil {
		r.Logout(ctx, player)
		return nil, nil
	}

	chosenRace := selectRace(ctx, session, r.db, player)
	if chosenRace == nil {
		r.Logout(ctx, player)
		return nil, nil
	}

	player.ColorProfile = *colorProfile
	player.HP = int32(chosenCharacterClass.HPAtFirstLevel)
	player.HPMax = int32(chosenCharacterClass.HPAtFirstLevel)
	player.Movement = 100
	player.MovementMax = 100
	player.UUID = uuid.New().String()
	player.Session = session
	player.CharacterClass = *chosenCharacterClass
	player.Race = *chosenRace

	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Exec(ctx, "INSERT INTO players (uuid, character_class, race, subrace, name, area, room, hp, hp_max, movement, movement_max, color_profile, password, logged_in) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		player.UUID, player.CharacterClass.ArchetypeSlug, player.Race.Slug, player.Race.SubRaceSlug, player.Name, player.AreaUUID, player.RoomUUID, player.HP, player.HPMax, player.Movement, player.MovementMax, player.ColorProfile.GetUUID(), player.Password, true)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to insert player: %v", err)
	}

	// TODO: everything below this line is probably junk.  Need to build player character based off choices made above.

	err = tx.Exec(ctx, "INSERT INTO player_abilities (uuid, player_uuid, strength, dexterity, constitution, intelligence, wisdom, charisma) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		uuid.New(), player.UUID, 10, 10, 10, 10, 10, 10)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to set player abilities: %v", err)
	}

	err = tx.Exec(ctx, "INSERT INTO player_equipments (uuid, player_uuid, Head, Neck, Chest, Arms, Hands, DominantHand, OffHand, Legs, Feet) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		uuid.New(), player.UUID, "", "", "", "", "", "", "", "", "", "")
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to set player equipments: %v", err)
	}
	player.Equipment = *NewPlayerEquipment()

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	player.Session = session

	return player, nil
}

func (r *Repository) LogoutAll(ctx context.Context) error {
	return r.db.Exec(ctx, "UPDATE players SET logged_in = 0")
}

func (r *Repository) Logout(ctx context.Context, player *Player) error {
	stmt, err := r.db.Prepare(ctx, "UPDATE players SET logged_in = FALSE WHERE uuid = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.Exec(ctx, player.UUID)
	if err != nil {
		return err
	}

	player.Session.Close()
	return nil
}

func (r *Repository) SetPlayerHealth(ctx context.Context, playerUUID string, health int) error {
	return r.db.Exec(ctx, "UPDATE players SET hp = ? WHERE uuid = ?", health, playerUUID)
}

func (r *Repository) GetColorProfileFromDB(ctx context.Context, uuid string) (*ColorProfile, error) {
	var colorProfile ColorProfile
	query := `SELECT uuid, name, primary_color, secondary_color, warning_color, danger_color, title_color, description_color FROM color_profiles WHERE uuid = ?`
	r.db.QueryRow(ctx, query, uuid).Scan(&colorProfile.UUID, &colorProfile.Name, &colorProfile.Primary, &colorProfile.Secondary, &colorProfile.Warning, &colorProfile.Danger, &colorProfile.Title, &colorProfile.Description)
	return &colorProfile, nil
}

func (r *Repository) GetPlayersInRoom(ctx context.Context, roomUUID string) ([]*Player, error) {
	var players []*Player
	query := `SELECT uuid, name FROM players WHERE room = ? and logged_in = 1`
	rows, err := r.db.Query(ctx, query, roomUUID)
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

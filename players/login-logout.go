package players

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

func getPlayerInput(reader io.Reader) string {
	r := bufio.NewReader(reader)
	input, _ := r.ReadString('\n')
	return strings.TrimSpace(input)
}

func getPlayerFromDB(db *sql.DB, playerName string) (*Player, error) {
	var player Player
	var colorProfile = &ColorProfile{}
	query := `SELECT p.name, p.uuid, p.area, p.room, p.health, p.health_max, p.movement, p.movement_max, p.mana, p.mana_max, p.password, cp.uuid, cp.name, cp.primary_color, cp.secondary_color, cp.warning_color, cp.danger_color, cp.title_color, cp.description_color
				FROM players p JOIN color_profiles cp ON cp.uuid = p.color_profile
				WHERE p.name = ?`
	err := db.QueryRow(query, playerName).
		Scan(&player.Name, &player.UUID, &player.Area, &player.Room, &player.Health, &player.HealthMax, &player.Movement, &player.MovementMax, &player.Mana, &player.ManaMax, &player.Password, &colorProfile.UUID, &colorProfile.Name, &colorProfile.Primary, &colorProfile.Secondary, &colorProfile.Warning, &colorProfile.Danger, &colorProfile.Title, &colorProfile.Description)
	if err != nil {
		return &player, err
	}

	player.ColorProfile = colorProfile

	return &player, nil
}

func setPlayerLoggedInStatus(db *sql.DB, playerUUID string, loggedIn bool) error {
	_, err := db.Exec("UPDATE players SET logged_in = ? WHERE uuid = ?", loggedIn, playerUUID)
	if err != nil {
		return err
	}
	return nil
}

func HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	return string(hashedPassword)
}

func getColorProfileFromDB(db *sql.DB, colorProfileUUID string) (*ColorProfile, error) {
	var colorProfile ColorProfile
	query := `SELECT uuid, name, primary_color, secondary_color, warning_color, danger_color, title_color, description_color
				FROM color_profiles
				WHERE uuid = ?`
	err := db.QueryRow(query, colorProfileUUID).
		Scan(&colorProfile.UUID, &colorProfile.Name, &colorProfile.Primary, &colorProfile.Secondary, &colorProfile.Warning, &colorProfile.Danger, &colorProfile.Title, &colorProfile.Description)
	if err != nil {
		return nil, err
	}

	return &colorProfile, nil
}

func createPlayer(conn net.Conn, db *sql.DB, playerName string) (*Player, error) {
	player := &Player{}
	player.Name = playerName

	fmt.Fprintf(conn, "Please enter a password you'd like to use: ")
	password := getPlayerInput(conn)
	player.Password = HashPassword(password)

	// default start point
	player.Area = "d71e8cf1-d5ba-426c-8915-4c7f5b22e3a9"
	player.Room = "189a729d-4e40-4184-a732-e2c45c66ff46"
	player.UUID = uuid.New().String()

	// "default" light mode color profile.  Should let the user choose?
	colorProfile, err := getColorProfileFromDB(db, "2c7dfd5b-d160-42e0-accb-b77d9686dbea")
	if err != nil {
		return nil, err
	}

	player.ColorProfile = colorProfile
	player.Health = 100
	player.HealthMax = 100
	player.Movement = 100
	player.MovementMax = 100
	player.Mana = 100
	player.ManaMax = 100
	player.UUID = uuid.New().String()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.Exec("INSERT INTO players (uuid, name, area, room, health, health_max, movement, movement_max, mana, mana_max, color_profile, password) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		player.UUID, player.Name, player.Area, player.Room, player.Health, player.HealthMax, player.Movement, player.MovementMax, player.Mana, player.ManaMax, player.ColorProfile.GetUUID(), player.Password)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to insert player: %v", err)
	}

	// Max out player's attributes while they're still a noob.
	_, err = tx.Exec("INSERT INTO player_attributes (uuid, player_uuid, strength, dexterity, constitution, intelligence, wisdom, charisma) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		uuid.New(), player.UUID, 18, 18, 18, 18, 18, 18)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to set player attributes: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	player.Conn = conn

	return player, nil
}

func LoginPlayer(conn net.Conn, db *sql.DB) (*Player, error) {

	fmt.Fprintf(conn, "Welcome! Please enter your player name: ")
	playerName := getPlayerInput(conn)

	player, err := getPlayerFromDB(db, playerName)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Fprintf(conn, "Player not found.  Do you want to create a new player? (y/n): ")
			answer := getPlayerInput(conn)

			if strings.ToLower(answer) == "y" {
				player, err = createPlayer(conn, db, playerName)
				if err != nil {
					return nil, err
				}
				return player, nil
			} else {
				return nil, errors.New("Player does not exist")
			}
		}
		return nil, err
	}

	fmt.Fprintf(conn, "Please enter your password: ")
	passwd := getPlayerInput(conn)
	err = bcrypt.CompareHashAndPassword([]byte(player.GetHashedPassword()), []byte(passwd))
	if err != nil {
		return nil, err
	}

	player.Conn = conn

	err = setPlayerLoggedInStatus(db, player.UUID, true)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func (player *Player) Logout(db *sql.DB) error {
	stmt, err := db.Prepare("UPDATE players SET logged_in = FALSE WHERE uuid = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(player.UUID)
	if err != nil {
		return err
	}

	player.Conn.Close()
	return nil
}

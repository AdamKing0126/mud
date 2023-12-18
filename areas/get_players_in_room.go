package areas

import (
	"database/sql"
	"fmt"
	"mud/interfaces"
	"net"
)

type ColorProfile struct {
}

type Player struct {
	UUID         string
	Name         string
	Area         string
	Health       int
	Mana         int
	Movement     int
	Room         string
	Conn         net.Conn
	Commands     []string
	ColorProfile interfaces.ColorProfileInterface
	LoggedIn     bool
	Password     string
}

func (player *Player) GetUUID() string {
	return player.UUID
}

func (player *Player) GetName() string {
	return player.Name
}

func (player *Player) GetRoom() string {
	return player.Room
}

func (player *Player) GetArea() string {
	return player.Area
}

func (player *Player) GetColorProfile() interfaces.ColorProfileInterface {
	return player.ColorProfile
}

func (player *Player) GetCommands() []string {
	return player.Commands
}

func (player *Player) GetConn() net.Conn {
	return player.Conn
}

func (player *Player) GetHealth() int {
	return player.Health
}

func (player *Player) GetMana() int {
	return player.Mana
}

func (player *Player) GetMovement() int {
	return player.Movement
}

func (player *Player) GetLoggedIn() bool {
	return player.LoggedIn
}

func GetPlayersInRoom(db *sql.DB, roomUUID string) ([]Player, error) {
	query := `
		SELECT uuid, name 
		FROM players 
		WHERE room = ? and logged_in = 1
	`
	rows, err := db.Query(query, roomUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var players []Player
	for rows.Next() {
		var player Player
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

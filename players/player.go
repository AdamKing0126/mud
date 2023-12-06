package players

import (
	"database/sql"
	"fmt"
	"mud/areas"
	"mud/display"
	"mud/interfaces"
	"net"
)

func NewPlayer(conn net.Conn) *Player {
	return &Player{Conn: conn}
}

type Player struct {
	UUID         string
	Name         string
	Room         string
	Area         string
	Health       int
	Conn         net.Conn
	Commands     []string
	ColorProfile interfaces.ColorProfileInterface
	LoggedIn     bool
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

func (player *Player) GetHealth() int {
	return player.Health
}

func (player *Player) GetConn() net.Conn {
	return player.Conn
}

func (player *Player) GetLoggedIn() bool {
	return player.LoggedIn
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

func (p *Player) GetCommands() []string {
	return p.Commands
}

func (p *Player) GetColorProfile() interfaces.ColorProfileInterface {
	return p.ColorProfile
}

func (p *Player) SetCommands(commands []string) {
	p.Commands = commands
}

func (player *Player) SetLocation(db *sql.DB, roomUUID string) error {
	area_rows, err := db.Query("SELECT area_uuid FROM rooms WHERE uuid=?", roomUUID)
	if err != nil {
		return fmt.Errorf("error retrieving area: %v", err)
	}
	defer area_rows.Close()

	if !area_rows.Next() {
		return fmt.Errorf("room with UUID %s does not have an area", roomUUID)
	}

	area := &areas.Area{}
	err = area_rows.Scan(&area.UUID)
	if err != nil {
		return err
	}

	player.Area = area.UUID
	player.Room = roomUUID

	area_rows.Close()

	stmt, err := db.Prepare("UPDATE players SET area = ?, room = ? WHERE uuid = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(area.UUID, roomUUID, player.UUID)
	if err != nil {
		return err
	}

	stmt.Close()
	return nil
}

type ColorProfile struct {
	UUID        string
	Name        string
	Primary     string
	Secondary   string
	Warning     string
	Danger      string
	Title       string
	Description string
}

func (c *ColorProfile) GetUUID() string {
	return c.UUID
}

func (c *ColorProfile) GetColor(colorUse string) string {
	switch colorUse {
	case "primary":
		return c.Primary
	case "secondary":
		return c.Secondary
	case "warning":
		return c.Warning
	case "danger":
		return c.Danger
	case "title":
		return c.Title
	default:
		return display.Reset
	}
}

func NewColorProfileFromDB(db *sql.DB, uuid string) (interfaces.ColorProfileInterface, error) {
	var colorProfile ColorProfile
	err := db.QueryRow("SELECT uuid, name, primary_color, secondary_color, warning_color, danger_color, title_color, description_color FROM color_profiles WHERE uuid = ?", uuid).
		Scan(&colorProfile.UUID, &colorProfile.Name, &colorProfile.Primary, &colorProfile.Secondary, &colorProfile.Warning, &colorProfile.Danger, &colorProfile.Title, &colorProfile.Description)
	if err != nil {
		return nil, err
	}
	return &colorProfile, nil
}

func GetPlayerByName(db *sql.DB, name string) (interfaces.PlayerInterface, error) {
	var player Player
	err := db.QueryRow("SELECT uuid, name, room, area, health, logged_in FROM players WHERE LOWER(name) = LOWER(?)", name).
		Scan(&player.UUID, &player.Name, &player.Room, &player.Area, &player.Health, &player.LoggedIn)
	if err != nil {
		return nil, err
	}
	return &player, nil
}

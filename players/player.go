package players

import (
	"database/sql"
	"fmt"
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
	HealthMax    int
	Mana         int
	ManaMax      int
	Movement     int
	MovementMax  int
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

func (player *Player) GetHealth() int {
	return player.Health
}

func (player *Player) GetHealthMax() int {
	return player.HealthMax
}

func (player *Player) GetMana() int {
	return player.Mana
}

func (player *Player) GetManaMax() int {
	return player.ManaMax
}

func (player *Player) GetMovement() int {
	return player.Movement
}

func (player *Player) GetMovementMax() int {
	return player.MovementMax
}

func (player *Player) GetHashedPassword() string {
	return player.Password
}

func (player *Player) SetHealth(health int) {
	player.Health = health
}

func (player *Player) SetMana(mana int) {
	player.Mana = mana
}

func (player *Player) SetMovement(movement int) {
	player.Movement = movement
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

	var areaUUID string
	err = area_rows.Scan(&areaUUID)
	if err != nil {
		return err
	}

	player.Area = areaUUID
	player.Room = roomUUID

	area_rows.Close()

	stmt, err := db.Prepare("UPDATE players SET area = ?, room = ?, movement = ? WHERE uuid = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	player.SetMovement(player.GetMovement() - 1)
	newMovement := player.GetMovement()
	_, err = stmt.Exec(areaUUID, roomUUID, newMovement, player.UUID)
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
	err := db.QueryRow("SELECT uuid, name, room, area, health, movement, mana, logged_in FROM players WHERE LOWER(name) = LOWER(?)", name).
		Scan(&player.UUID, &player.Name, &player.Room, &player.Area, &player.Health, &player.Movement, &player.Mana, &player.LoggedIn)
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func calculateHealthRegen(p *Player) float64 {
	return 1.1
}

func calculateManaRegen(p *Player) float64 {
	return 1.1
}

func calculateMovementRegen(p *Player) float64 {
	return 1.1
}

func (p *Player) Regen(db *sql.DB) error {
	healthRegen := calculateHealthRegen(p)
	manaRegen := calculateManaRegen(p)
	movementRegen := calculateMovementRegen(p)

	p.Health = int(float64(p.Health) * healthRegen)
	if p.Health > p.HealthMax {
		p.Health = p.HealthMax
	}

	p.Mana = int(float64(p.Mana) * manaRegen)
	if p.Mana > p.ManaMax {
		p.Mana = p.ManaMax
	}

	p.Movement = int(float64(p.Movement) * movementRegen)
	if p.Movement > p.MovementMax {
		p.Movement = p.MovementMax
	}

	_, err := db.Exec("UPDATE players SET health = ?, mana = ?, movement = ? WHERE uuid = ?", p.Health, p.Mana, p.Movement, p.UUID)
	if err != nil {
		return err
	}
	return nil
}

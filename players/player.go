package players

import (
	"database/sql"
	"fmt"
	"mud/interfaces"
	"net"
	"reflect"
)

func NewPlayer(conn net.Conn) *Player {
	return &Player{Conn: conn}
}

type Player struct {
	UUID            string
	Name            string
	Room            string
	Area            string
	Health          int
	HealthMax       int
	Mana            int
	ManaMax         int
	Movement        int
	MovementMax     int
	Conn            net.Conn
	Commands        []string
	ColorProfile    interfaces.ColorProfileInterface
	LoggedIn        bool
	Password        string
	PlayerAbilities interfaces.PlayerAbilitiesInterface
	Equipment       PlayerEquipment
}

func (player *Player) GetArmorClass() int {
	// 10 + armor_bonus + shield_bonus + dexterity_modifier + other_modifiers
	base := 10
	armorBonus := 0
	shieldBonus := 0
	dexModifier := player.PlayerAbilities.GetDexterityModifier()
	otherModifiers := 0
	return base + armorBonus + shieldBonus + dexModifier + otherModifiers
}

func (player *Player) GetSizeModifier() int {
	// Need to update this.  Probably need to move this out, so it can be used by players and monsters

	sizeTable := map[string]int{
		"colossal":   -8,
		"gargantuan": -4,
		"huge":       -2,
		"large":      -1,
		"medium":     0,
		"small":      1,
		"tiny":       2,
		"diminutive": 4,
		"fine":       8,
	}
	return sizeTable["medium"]
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

func (player *Player) GetCommands() []string {
	return player.Commands
}

func (player *Player) GetColorProfile() interfaces.ColorProfileInterface {
	return player.ColorProfile
}

func (player *Player) GetAbilities() interfaces.AbilitiesInterface {
	return player.PlayerAbilities
}

func (player *Player) SetCommands(commands []string) {
	player.Commands = commands
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

func (player *Player) SetAbilities(abilities interfaces.PlayerAbilitiesInterface) {
	player.PlayerAbilities = abilities
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

func (player *Player) Equip(db *sql.DB, item interfaces.ItemInterface) bool {
	// get the location where the thing goes
	val := reflect.ValueOf(&player.Equipment).Elem()
	itemEquipSlots := []string{}
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)

		for _, slot := range item.GetEquipmentSlots() {
			if string(slot) == field.Name {
				itemEquipSlots = append(itemEquipSlots, field.Name)
			}
		}

	}
	fmt.Println(itemEquipSlots)

	// generate the query to retrieve columns from player_equipments table
	queryString := "SELECT "
	for i, slot := range itemEquipSlots {
		if i > 0 {
			queryString += ", "
		}
		queryString += slot
	}
	queryString += " FROM player_equipments WHERE player_uuid = ? LIMIT 1"
	rows, err := db.Query(queryString, player.GetUUID())
	if err != nil {
		fmt.Printf("error retrieving player equipments: %v", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		fmt.Printf("error getting columns: %v", err)
	}

	values := make([]interface{}, len(columns))
	pointers := make([]interface{}, len(columns))

	for i := range columns {
		pointers[i] = &values[i]
	}

	if rows.Next() {
		err = rows.Scan(pointers...)
		if err != nil {
			fmt.Printf("error scanning row: %v", err)
			return false
		}
	}

	// iterate through the values
	// get the index of the first empty value, equip the item there and then break
	for idx, val := range values {
		if val == "" {
			itemUUID := item.GetUUID()
			queryString := "UPDATE player_equipments SET "
			queryString += columns[idx]
			queryString += " = ? WHERE player_uuid = ?"
			rows.Close()
			_, err = db.Exec(queryString, itemUUID, player.GetUUID())
			if err != nil {
				fmt.Printf("error inserting into player_equipments: %v", err)
				return false
			}

			return true
		}

	}

	// if no empty values, remove the first item with a value,
	fmt.Printf("equipping at first slot.")
	// then equip the new item

	fmt.Printf("yo, dude %v", columns)
	return true

}

func GetPlayer(db *sql.DB, playerName string) (interfaces.PlayerInterface, error) {
	var player Player
	err := db.QueryRow("SELECT uuid, name, room, area, health, movement, mana, logged_in, FROM players WHERE LOWER(p.name) = LOWER(?)", playerName).
		Scan(&player.UUID, &player.Name, &player.Room, &player.Area, &player.Health, &player.Movement, &player.Mana, &player.LoggedIn)
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func GetPlayerByName(db *sql.DB, name string) (interfaces.PlayerInterface, error) {
	var player Player
	var playerAbilities PlayerAbilities
	err := db.QueryRow("SELECT p.uuid, p.name, p.room, p.area, p.health, p.movement, p.mana, p.logged_in, pa.intelligence, pa.dexterity, pa.charisma, pa.constitution, pa.wisdom, pa.strength FROM players p JOIN player_attributes pa ON p.uuid = pa.player_uuid WHERE LOWER(p.name) = LOWER(?)", name).
		Scan(&player.UUID, &player.Name, &player.Room, &player.Area, &player.Health, &player.Movement, &player.Mana, &player.LoggedIn, &playerAbilities.Intelligence, &playerAbilities.Dexterity, &playerAbilities.Charisma, &playerAbilities.Constitution, &playerAbilities.Wisdom, &playerAbilities.Strength)
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

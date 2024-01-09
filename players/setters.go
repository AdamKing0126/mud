package players

import (
	"database/sql"
	"fmt"
	"mud/interfaces"
	"net"
)

func (player *Player) SetAbilities(abilities interfaces.PlayerAbilities) {
	playerAbilities, ok := abilities.(*PlayerAbilities)
	if !ok {
		fmt.Errorf("error setting abilities")
	}
	player.PlayerAbilities = *playerAbilities
}

func (player *Player) SetCommands(commands []string) {
	player.Commands = commands
}

func (player *Player) SetConn(conn net.Conn) {
	player.Conn = conn
}

func (player *Player) SetHealth(health int) {
	player.Health = health
}

func (player *Player) SetInventory(inventory []interfaces.Item) {
	player.Inventory = inventory
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

func (player *Player) SetMana(mana int) {
	player.Mana = mana
}

func (player *Player) SetMovement(movement int) {
	player.Movement = movement
}

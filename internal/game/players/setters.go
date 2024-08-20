package players

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func (player *Player) SetLocation(db *sqlx.DB, roomUUID string) error {
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

	player.AreaUUID = areaUUID
	player.RoomUUID = roomUUID

	area_rows.Close()

	stmt, err := db.Prepare("UPDATE players SET area = ?, room = ?, movement = ? WHERE uuid = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	player.Movement--
	_, err = stmt.Exec(areaUUID, roomUUID, player.Movement, player.UUID)
	if err != nil {
		return err
	}

	stmt.Close()
	return nil
}

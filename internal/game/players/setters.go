package players

import (
	"context"
	"fmt"

	"github.com/adamking0126/mud/pkg/database"
)

func (player *Player) SetLocation(ctx context.Context, db database.DB, roomUUID string) error {
	area_rows, err := db.Query(ctx, "SELECT area_uuid FROM rooms WHERE uuid=?", roomUUID)
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

	stmt, err := db.Prepare(ctx, "UPDATE players SET area = ?, room = ?, movement = ? WHERE uuid = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	player.Movement--
	err = stmt.Exec(ctx, areaUUID, roomUUID, player.Movement, player.UUID)
	if err != nil {
		return err
	}

	stmt.Close()
	return nil
}

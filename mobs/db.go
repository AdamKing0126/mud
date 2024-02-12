package mobs

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/mitchellh/mapstructure"
)

func GetMobsInRoom(db *sqlx.DB, roomUUID string) ([]Mob, error) {
	var mobs []Mob
	rows, err := db.Queryx("SELECT * FROM mobs WHERE room_uuid = ?", roomUUID)
	if err != nil {
		log.Fatalf("failed to fetch mobs from tabel: %v", err)
	}

	for rows.Next() {
		result := make(map[string]interface{})
		err = rows.MapScan(result)
		if err != nil {
			log.Fatalf("failed to scan row: %v", err)
		}

		var mob Mob
		err = mapstructure.Decode(result, &mob)
		if err != nil {
			log.Fatalf("failed to decode row: %v", err)
		}
		mobs = append(mobs, mob)
	}

	return mobs, nil

}

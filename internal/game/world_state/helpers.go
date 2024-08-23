package world_state

import (
	"context"
	"fmt"

	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/internal/game/items"
	"github.com/adamking0126/mud/internal/game/mobs"
	"github.com/adamking0126/mud/pkg/database"
)

// todo - need to migrate this over to areas service/repository
func setMobs(ctx context.Context, db database.DB, roomInArea *areas.Room) {
	// retrieve the mobs and attach them to the room in the WorldState
	mobsInRoom, err := mobs.GetMobsInRoom(ctx, db, roomInArea.UUID)
	if err != nil {
		fmt.Printf("error retrieving mobs: %v", err)
	}
	roomInArea.Mobs = mobsInRoom
}

// todo - need to migrate this over to areas service/repository
func setItems(ctx context.Context, db database.DB, roomInArea *areas.Room) {
	// retrieve the items and attach them to the room in the WorldState
	itemsInRoom, err := items.GetItemsInRoom(ctx, db, roomInArea.UUID)
	if err != nil {
		fmt.Printf("error retrieving items: %v", err)
	}
	roomInArea.Items = itemsInRoom
}

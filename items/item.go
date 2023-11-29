package items

import (
	"database/sql"
	"fmt"

	"mud/interfaces"

	"github.com/google/uuid"
)

type ItemLocation struct {
	ItemUUID   uuid.UUID
	RoomUUID   uuid.UUID
	PlayerUUID uuid.UUID
}

type Item struct {
	UUID        string
	Name        string
	Description string
}

func (item *Item) GetUUID() string {
	return item.UUID
}

func (item *Item) GetName() string {
	return item.Name
}

func (item *Item) GetDescription() string {
	return item.Description
}

func GetItemsForPlayer(db *sql.DB, playerUUID string) ([]interfaces.ItemInterface, error) {
	query := `
		SELECT i.uuid, i.name, i.description
		FROM item_locations il
		JOIN items i on il.item_uuid = i.uuid
		WHERE il.player_uuid = ?
	`

	rows, err := db.Query(query, playerUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.UUID, &item.Name, &item.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	itemInterfaces := make([]interfaces.ItemInterface, len(items))
	for i, item := range items {
		itemInterfaces[i] = &item
	}
	return itemInterfaces, nil
}

func GetItemsInRoom(db *sql.DB, roomUUID string) ([]interfaces.ItemInterface, error) {
	query := `
		SELECT i.uuid, i.name, i.description
		FROM item_locations il
		JOIN items i ON il.item_uuid = i.uuid
		WHERE il.room_uuid = ?
	`
	rows, err := db.Query(query, roomUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.UUID, &item.Name, &item.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	itemInterfaces := make([]interfaces.ItemInterface, len(items))
	for i := range items {
		itemInterfaces[i] = &items[i]
	}

	return itemInterfaces, nil
}

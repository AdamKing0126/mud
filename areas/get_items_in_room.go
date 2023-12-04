package areas

import (
	"database/sql"
	"fmt"
)

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

func GetItemsInRoom(db *sql.DB, roomUUID string) ([]Item, error) {
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

	return items, nil
}

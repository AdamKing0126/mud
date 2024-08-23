package world_state

import (
	"context"

	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/pkg/database"
)

type Repository struct {
	db database.DB
}

func NewRepository(db database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetRoom(ctx context.Context, roomUUID string, includeExits bool) *areas.Room {
	return nil
}

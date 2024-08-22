package areas

import (
	"context"

	"github.com/adamking0126/mud/internal/game/items"
	"github.com/adamking0126/mud/pkg/database"
)

type Service struct {
	repo *Repository
}

func NewService(db database.DB) *Service {
	return &Service{repo: NewRepository(db)}
}

func (s *Service) GetRoom(ctx context.Context, roomUUID string) (*Room, error) {
	return s.repo.GetRoomFromDB(ctx, roomUUID)
}

func (s *Service) AddItemToRoom(ctx context.Context, room *Room, item *items.Item) error {
	return s.repo.AddItemToRoom(ctx, room, item)
}

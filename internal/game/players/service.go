package players

import (
	"context"

	"github.com/adamking0126/mud/pkg/database"
)

type Service struct {
	repo *Repository
}

func NewService(db database.DB) *Service {
	return &Service{
		repo: NewRepository(db),
	}
}

func (s *Service) GetPlayerByName(ctx context.Context, name string) (*Player, error) {
	return s.repo.GetPlayerByName(ctx, name)
}

func (s *Service) GetPlayerByNameFull(ctx context.Context, name string) (*Player, error) {
	return s.repo.GetPlayerByNameFull(ctx, name)
}

func (s *Service) SetPlayerLoggedInStatus(ctx context.Context, playerUUID string, loggedIn bool) error {
	return s.repo.SetPlayerLoggedInStatus(ctx, playerUUID, loggedIn)
}

// Add other methods that use the repository...

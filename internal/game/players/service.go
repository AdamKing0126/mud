package players

import (
	"context"

	"github.com/adamking0126/mud/internal/game/items"
	"github.com/adamking0126/mud/pkg/database"
	"github.com/charmbracelet/ssh"
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

func (s *Service) GetPlayerFromDB(ctx context.Context, name string) (*Player, error) {
	return s.repo.GetPlayerFromDB(ctx, name)
}

func (s *Service) CreatePlayer(ctx context.Context, session ssh.Session, playerName string) (*Player, error) {
	return s.repo.CreatePlayer(ctx, session, playerName)
}

func (s *Service) GetColorProfileForPlayer(ctx context.Context, playerUUID string) *ColorProfile {
	return s.repo.GetColorProfileForPlayerByUUID(ctx, playerUUID)
}

func (s *Service) SetPlayerColorProfile(ctx context.Context, player *Player) error {
	colorProfile := s.GetColorProfileForPlayer(ctx, player.UUID)
	player.ColorProfile = *colorProfile
	return nil
}

func (s *Service) SetPlayerEquipment(ctx context.Context, player *Player) error {
	equipment := s.repo.GetEquipmentForPlayerByUUID(ctx, player.UUID)
	player.Equipment = *equipment
	return nil
}

func (s *Service) SetPlayerInventory(ctx context.Context, player *Player) error {
	inventory := s.repo.GetInventoryForPlayerByUUID(ctx, player.UUID)
	player.Inventory = inventory
	return nil
}

func (s *Service) SetPlayerLoggedInStatus(ctx context.Context, player *Player, loggedIn bool) error {
	return s.repo.SetPlayerLoggedInStatus(ctx, player.UUID, loggedIn)
}

func (s *Service) AddItemToPlayer(ctx context.Context, player *Player, item *items.Item) error {
	err := item.SetLocation(ctx, s.repo.db, player.UUID, "")
	if err != nil {
		return err
	}

	player.Inventory = append(player.Inventory, item)
	return nil
}

func (s *Service) SetPlayerHealth(ctx context.Context, player *Player, health int) error {
	return s.repo.SetPlayerHealth(ctx, player.UUID, health)
}

func (s *Service) LogoutAllPlayers(ctx context.Context) error {
	return s.repo.LogoutAll(ctx)
}

func (s *Service) LogoutPlayer(ctx context.Context, player *Player) error {
	return s.repo.Logout(ctx, player)
}

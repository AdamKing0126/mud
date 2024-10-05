package world_state

import (
	"context"

	"github.com/adamking0126/mud/internal/game/areas"
	"github.com/adamking0126/mud/internal/game/items"
	"github.com/adamking0126/mud/internal/game/players"
	"github.com/adamking0126/mud/pkg/database"
)

type Service struct {
	repo        *Repository
	areaService *areas.Service
}

func NewService(db database.DB, areaService *areas.Service) *Service {
	return &Service{repo: NewRepository(db), areaService: areaService}
}

func (s *Service) GetRoom(ctx context.Context, roomUUID string, includeExits bool) *areas.Room {
	roomUUIDToAreaMap := s.areaService.GetRoomUUIDToAreaMap()
	areaUUID := roomUUIDToAreaMap[roomUUID]

	roomCount := s.areaService.GetRoomCount(ctx, areaUUID)

	area := s.areaService.GetAreaUUIDToAreaMap()[areaUUID]
	if len(area.Rooms) == roomCount {
		for idx := range area.Rooms {
			if area.Rooms[idx].UUID == roomUUID {
				return area.Rooms[idx]
			}
		}
	}

	return s.repo.GetRoom(ctx, roomUUID, includeExits)
}

func (s *Service) GetRoomByUUID(ctx context.Context, roomUUID string) *areas.Room {
	return s.areaService.GetRoomByUUID(ctx, roomUUID)
}

func (s *Service) GetAreaByUUID(ctx context.Context, areaUUID string) *areas.Area {
	return s.areaService.GetAreaByUUID(ctx, areaUUID)
}

func (s *Service) RemovePlayerFromRoom(ctx context.Context, roomUUID string, player *players.Player) error {
	err := s.areaService.RemovePlayerFromRoom(ctx, roomUUID, player)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) AddPlayerToRoom(ctx context.Context, roomUUID string, player *players.Player) error {
	err := s.areaService.AddPlayerToRoom(ctx, roomUUID, player)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) AddItemToRoom(ctx context.Context, room *areas.Room, item *items.Item) {
	s.areaService.AddItemToRoom(ctx, room, item)
}

func (s *Service) RemoveItemFromRoom(ctx context.Context, room *areas.Room, item *items.Item) {
	s.areaService.RemoveItemFromRoom(ctx, room, item)
}
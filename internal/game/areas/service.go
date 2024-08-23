package areas

import (
	"context"
	"time"

	"github.com/adamking0126/mud/internal/game/items"
	"github.com/adamking0126/mud/internal/game/players"
	"github.com/adamking0126/mud/pkg/database"
)

type Service struct {
	repo              *Repository
	playerService     *players.Service
	areaUUIDToAreaMap map[string]*Area
	roomUUIDToAreaMap map[string]string
	areaChannels      map[string]chan Action
	connections       map[string]*players.Player
}

func NewService(db database.DB, playerService *players.Service, connections map[string]*players.Player) *Service {
	return &Service{repo: NewRepository(db), playerService: playerService, connections: connections}
}

func (s *Service) GetRoom(ctx context.Context, targetRoomUUID string, followExits bool) *Room {
	// By default, we do not load every room into memory.
	// We only load rooms that are in the same area as the player.
	// We do this by looking up the area of the room, and then querying the database for the number of rooms in that area.
	// If the number of rooms in the area is equal to the number of rooms that we have loaded into memory,
	// then we have all the rooms in the area loaded into memory.
	// If not, we need to load the rooms from the database.

	// There is currently not a process to "unload" rooms from memory, but in the future we may want to do that.

	area := s.areaUUIDToAreaMap[s.roomUUIDToAreaMap[targetRoomUUID]] // todo this seems janky

	roomCount, err := s.repo.GetRoomCount(ctx, area.UUID)
	if err != nil {
		return nil
	}
	if len(area.Rooms) == roomCount {
		// We have all the rooms in the area loaded into memory.
		for _, room := range area.Rooms {
			if room.UUID == targetRoomUUID {
				return room
			}
		}
	}

	retrievedRoom := s.repo.GetRoomFromDBAndLoadArea(ctx, area, targetRoomUUID)

	// now that we have loaded all the rooms in the area, we can go back and hook up all the
	// exits.  if one of the exits happens to exist in a different area, we can make a db query to retrieve that one.
	for _, room := range area.Rooms {
		if followExits {
			s.setExits(ctx, area, room)
		}
		s.setPlayers(ctx, room)

		// TODO - still need to add services for items and mobs
		// s.setItems(ctx, area, room)
		// s.setMobs(ctx, area, room)
	}
	return retrievedRoom
}

func (s *Service) GetRoomByUUID(ctx context.Context, roomUUID string) *Room {
	areaUUID := s.roomUUIDToAreaMap[roomUUID]
	area := s.areaUUIDToAreaMap[areaUUID]
	return s.GetRoomFromAreaByUUID(ctx, area, roomUUID)
}

func (s *Service) GetAreaByUUID(ctx context.Context, areaUUID string) *Area {
	return s.areaUUIDToAreaMap[areaUUID]
}

func (s *Service) AddItemToRoom(ctx context.Context, room *Room, item *items.Item) error {
	return s.repo.AddItemToRoom(ctx, room, item)
}

func (s *Service) RunAreaActions(ctx context.Context, areaUUID string, ch chan Action) {
	ticker := time.NewTicker(time.Second)
	tickerCounter := 0
	defer ticker.Stop()

	playerActionsMap := make(map[string]PlayerActions)

	for {
		select {
		case action := <-ch:
			player := action.GetPlayer()
			pa := playerActionsMap[player.UUID]
			pa.Actions = append(pa.Actions, action)
			playerActionsMap[player.UUID] = pa
		case <-ticker.C:
			tickerCounter++
			if tickerCounter%15 == 0 {
				s.processAreaBeat(ctx, areaUUID)
			}

			s.processPlayerActions(ctx, playerActionsMap)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) GetRoomCount(ctx context.Context, areaUUID string) int {
	roomCount, err := s.repo.GetRoomCount(ctx, areaUUID)
	if err != nil {
		return 0
	}
	return roomCount
}

func (s *Service) processAreaBeat(ctx context.Context, areaUUID string) {
	// Implement area beat logic here
	// You may need to add methods to the repository to fetch necessary data
}

func (s *Service) processPlayerActions(ctx context.Context, playerActionsMap map[string]PlayerActions) {
	// Implement player action processing logic here
	// You may need to add methods to the repository to update player data
}

func (s *Service) LoadAreas(ctx context.Context) error {
	areaRoomData, err := s.repo.GetAreasFromDB(ctx)
	if err != nil {
		return err
	}

	areaUUIDToAreaMap := make(map[string]*Area)
	roomUUIDToAreaMap := make(map[string]string)
	areaChannels := make(map[string]chan Action)

	for _, data := range areaRoomData {
		if _, ok := areaUUIDToAreaMap[data.AreaUUID]; !ok {
			areaUUIDToAreaMap[data.AreaUUID] = NewArea(data.AreaUUID, data.AreaName, data.AreaDescription)
			areaChannels[data.AreaUUID] = make(chan Action)
			go areaUUIDToAreaMap[data.AreaUUID].Run(ctx, s.repo.db, areaChannels[data.AreaUUID], s.connections)
		}
		roomUUIDToAreaMap[data.RoomUUID] = data.AreaUUID
	}

	s.areaUUIDToAreaMap = areaUUIDToAreaMap
	s.roomUUIDToAreaMap = roomUUIDToAreaMap
	s.areaChannels = areaChannels

	return nil
}

func (s *Service) GetAreaUUIDToAreaMap() map[string]*Area {
	return s.areaUUIDToAreaMap
}

func (s *Service) GetRoomUUIDToAreaMap() map[string]string {
	return s.roomUUIDToAreaMap
}

func (s *Service) GetAreaChannels() map[string]chan Action {
	return s.areaChannels
}

func (s *Service) GetRoomFromAreaByUUID(ctx context.Context, area *Area, roomUUID string) *Room {
	for _, room := range area.Rooms {
		if room.UUID == roomUUID {
			return room
		}
	}
	return nil
}

func (s *Service) setExits(ctx context.Context, area *Area, room *Room) {
	exits := room.Exits
	exitInfo := ExitInfo{}
	exitInfo.South = s.getExitRoom(ctx, area, exits.South)
	exitInfo.North = s.getExitRoom(ctx, area, exits.North)
	exitInfo.East = s.getExitRoom(ctx, area, exits.East)
	exitInfo.West = s.getExitRoom(ctx, area, exits.West)
	exitInfo.Up = s.getExitRoom(ctx, area, exits.Up)
	exitInfo.Down = s.getExitRoom(ctx, area, exits.Down)
	room.Exits = &exitInfo
}

func (s *Service) getExitRoom(ctx context.Context, area *Area, room *Room) *Room {
	if room != nil {
		if room.Name == "" {
			exitRoom := s.GetRoomFromAreaByUUID(ctx, area, room.UUID)
			if exitRoom != nil {
				return exitRoom
			}
			return s.repo.GetRoomByUUIDFromDB(ctx, room.UUID)
		}
	}
	return nil
}

func (s *Service) setPlayers(ctx context.Context, room *Room) {
	players := s.playerService.GetPlayersInRoom(ctx, room.UUID)
	room.Players = players
}

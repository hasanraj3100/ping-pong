package game

import (
	"hasanraj3100/ping-pong/internal/domain"
)

type GameService struct {
	repo        GameRepo
	broadcaster Broadcaster
}

func NewGameService(repo GameRepo, broadcaster Broadcaster) *GameService {
	return &GameService{repo: repo, broadcaster: broadcaster}
}

func (s *GameService) CreateGame(creator domain.User) domain.Game {
	g := s.repo.Create(creator)
	s.broadcaster.Broadcast(g.ID, g)
	return g
}

func (s *GameService) GetGame(gameID string) (domain.Game, bool) {
	return s.repo.FindByID(gameID)
}

func (s *GameService) JoinGame(gameID string, joiner domain.User) (domain.Game, error) {
	g, ok := s.repo.FindByID(gameID)
	if !ok {
		return domain.Game{}, ErrGameNotFound
	}

	if g.Player1.UUID == joiner.UUID {
		return domain.Game{}, ErrAlreadyJoined
	}

	if g.Player2.UUID != "" {
		return domain.Game{}, ErrGameFull
	}

	g.Player2 = joiner
	updated := s.repo.Update(g)
	s.broadcaster.Broadcast(updated.ID, updated)
	return updated, nil
}

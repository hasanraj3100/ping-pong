package game

import (
	"hasanraj3100/ping-pong/internal/domain"
)

type GameService struct {
	repo GameRepo
}

func NewGameService(repo GameRepo) *GameService {
	return &GameService{repo: repo}
}

func (s *GameService) CreateGame(creator domain.User) domain.Game {
	return s.repo.Create(creator)
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
	return s.repo.Update(g), nil
}

package repository

import (
	"sync"

	"hasanraj3100/ping-pong/internal/domain"

	"github.com/google/uuid"
)

type GameRepository struct {
	mu    sync.Mutex
	games map[string]domain.Game
}

func NewGameRepository() *GameRepository {
	return &GameRepository{
		games: make(map[string]domain.Game),
	}
}

func (repo *GameRepository) Create(creator domain.User) domain.Game {
	game := domain.Game{
		ID:      uuid.NewString(),
		Player1: creator,
		State: domain.GameData{
			Player1YPosition: 110,
			Player2YPosition: 110,
			Player1Ready:     false,
			Player2Ready:     false,
		},
	}

	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.games[game.ID] = game
	return game
}

func (repo *GameRepository) FindByID(id string) (domain.Game, bool) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	game, ok := repo.games[id]
	return game, ok
}

func (repo *GameRepository) Update(game domain.Game) domain.Game {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.games[game.ID] = game
	return game
}

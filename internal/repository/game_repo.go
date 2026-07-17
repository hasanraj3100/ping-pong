package repository

import (
	"hasanraj3100/ping-pong/internal/domain"

	"github.com/google/uuid"
)

type GameRepository struct {
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
		},
	}

	repo.games[game.ID] = game
	return game
}

func (repo *GameRepository) FindByID(id string) (domain.Game, bool) {
	game, ok := repo.games[id]
	return game, ok
}

func (repo *GameRepository) Update(game domain.Game) domain.Game {
	repo.games[game.ID] = game
	return game
}

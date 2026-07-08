package game

import "hasanraj3100/ping-pong/internal/domain"

type GameRepo interface {
	Create(c domain.User) domain.Game
	FindByID(id string) (domain.Game, bool)
	Update(g domain.Game) domain.Game
}

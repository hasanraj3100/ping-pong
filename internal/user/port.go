package user

import "hasanraj3100/ping-pong/internal/domain"

type Service interface {
	Register(name string) domain.User
}

type UserRepo interface {
	Create(name string) domain.User
}

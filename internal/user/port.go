package user

import "hasanraj3100/ping-pong/internal/domain"

type Service interface {
	Register(name string) domain.User
	GetByUUID(uuid string) (domain.User, bool)
}

type UserRepo interface {
	Create(name string) domain.User
	FindByUUID(uuid string) (domain.User, bool)
}

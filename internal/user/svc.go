package user

import "hasanraj3100/ping-pong/internal/domain"

type userService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) Service {
	return &userService{
		repo: repo,
	}
}

func (s *userService) Register(name string) domain.User {
	return s.repo.Create(name)
}

func (s *userService) GetByUUID(uuid string) (domain.User, bool) {
	return s.repo.FindByUUID(uuid)
}

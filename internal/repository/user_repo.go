package repository

import (
	"log"

	"hasanraj3100/ping-pong/internal/domain"

	"github.com/google/uuid"
)

type userRepository struct {
	users map[string]domain.User
}

func NewUserRepository() *userRepository {
	return &userRepository{
		users: make(map[string]domain.User),
	}
}

func (us *userRepository) Create(name string) domain.User {
	user := domain.User{
		UUID: uuid.NewString(),
		Name: name,
	}

	us.users[user.UUID] = user
	log.Println("Created User", user.Name, user.UUID)

	return user
}

func (us *userRepository) FindByUUID(uuid string) (domain.User, bool) {
	user, ok := us.users[uuid]
	return user, ok
}

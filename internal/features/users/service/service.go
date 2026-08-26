package users_service

import (
	"context"

	"github.com/google/uuid"

	"github.com/kupr666/Orange_Team/internal/core/domain"
)

type UsersService struct {
	usersRepository UsersRepository
}

type UsersRepository interface {
	GetUser(
		ctx context.Context,
		userID uuid.UUID,
	) (domain.User, error)

	PatchUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	DeleteUser(
		ctx context.Context,
		userID uuid.UUID,
	) error
}

func NewUsersService(repo UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: repo,
	}
}

package auth_service

import (
	"context"

	"github.com/kupr666/Orange_Team/internal/core/domain"
)

type AuthService struct {
	authRepository AuthRepository
}

func NewAuthService(repo AuthRepository) *AuthService {
	return &AuthService{
		authRepository: repo,
	}
}

type AuthRepository interface {
	RegisterUser(
		ctx context.Context,
		email string,
		passwordHash string,
		fullName string,
	) (domain.User, error)
}

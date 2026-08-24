package authentication_service

import (
	"context"

	"github.com/kupr666/Orange_Team/internal/core/domain"
)

type AuthenticationService struct {
	authenticationRepository AuthenticationRepository
}

func NewAuthenticationService(repo AuthenticationRepository) *AuthenticationService {
	return &AuthenticationService{
		authenticationRepository: repo,
	}
}

type AuthenticationRepository interface {
	CreateUser(
		ctx context.Context,
		email string,
		passwordHash string,
		fullName string,
	) (domain.User, error)
}

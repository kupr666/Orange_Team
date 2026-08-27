package authentication_service

import (
	"context"
	"fmt"

	core_auth "github.com/kupr666/Orange_Team/internal/core/auth"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	authentication_domain "github.com/kupr666/Orange_Team/internal/features/authentication/domain"
	"golang.org/x/crypto/bcrypt"
)

type AccessTokenIssuer interface {
	IssueAccessToken(
		principal core_auth.Principal,
	) (string, error)
}

type AuthenticationService struct {
	authenticationRepository AuthenticationRepository
	accessTokenIssuer        AccessTokenIssuer
	dummyPasswordHash        []byte
}

type AuthenticationRepository interface {
	RegisterUser(
		ctx context.Context,
		email string,
		passwordHash string,
		fullName string,
	) (domain.User, error)
	GetLoginCredentialsByEmail(
		ctx context.Context,
		email string,
	) (authentication_domain.StoredCredentials, error)
}

func NewAuthenticationService(
	repo AuthenticationRepository,
	accessTokenIssuer AccessTokenIssuer,
) (*AuthenticationService, error) {
	dummyPasswordHash, err := bcrypt.GenerateFromPassword(
		[]byte("dummy-pass-used-only-for-timing"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, fmt.Errorf("generate dummy password hash: %w", err)
	}
	return &AuthenticationService{
		authenticationRepository: repo,
		accessTokenIssuer:        accessTokenIssuer,
		dummyPasswordHash:        dummyPasswordHash,
	}, nil
}

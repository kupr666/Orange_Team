package authentication_service

import (
	"context"
	"errors"
	"fmt"

	core_auth "github.com/kupr666/Orange_Team/internal/core/auth"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	"golang.org/x/crypto/bcrypt"
)

const tokenTypeBearer = "Bearer"

var errInvalidCredentials = fmt.Errorf(
	"invalid login credentials: %w",
	core_errors.ErrUnauthorized,
)

type LoginResult struct {
	AccessToken string
	TokenType   string
}

func (s *AuthenticationService) Login(
	ctx context.Context,
	email string,
	password string,
) (LoginResult, error) {

	email = normalizeEmail(email)

	if err := validateEmail(email); err != nil {
		return LoginResult{}, err
	}

	if err := validatePassword(password); err != nil {
		return LoginResult{}, err
	}

	credentials, err := s.authenticationRepository.GetLoginCredentialsByEmail(ctx, email)
	userFound := true

	if err != nil {
		if !errors.Is(err, core_errors.ErrNotFound) {
			return LoginResult{}, fmt.Errorf(
				"get stored credentials: %w",
				err,
			)
		}

		userFound = false
		credentials.PasswordHash = string(s.dummyPasswordHash)
	}

	compareErr := bcrypt.CompareHashAndPassword(
		[]byte(credentials.PasswordHash),
		[]byte(password),
	)

	if !userFound || errors.Is(compareErr, bcrypt.ErrMismatchedHashAndPassword) {
		return LoginResult{}, errInvalidCredentials
	}
	if compareErr != nil {
		return LoginResult{}, fmt.Errorf("copare password hash: %w", compareErr)
	}

	principal := core_auth.Principal{
		UserID: credentials.UserID,
		Role:   credentials.Role,
	}

	accessToken, err := s.accessTokenIssuer.IssueAccessToken(principal)
	if err != nil {
		return LoginResult{}, fmt.Errorf(
			"issue access token: %w",
			err,
		)
	}

	return LoginResult{
		AccessToken: accessToken,
		TokenType:   tokenTypeBearer,
	}, nil
}

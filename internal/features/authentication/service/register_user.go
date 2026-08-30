package authentication_service

import (
	"context"
	"fmt"
	"regexp"

	"github.com/kupr666/Orange_Team/internal/core/domain"
	"golang.org/x/crypto/bcrypt"
)

const (
	minimumPasswordLengthBytes = 8
	maximumPasswordLengthBytes = 72
	minimumEmailLength         = 5
	maximumEmailLength         = 30
	minimumFullNameLength      = 2
	maximumFullNameLength      = 50
)

var emailPattern = regexp.MustCompile(
	`^[a-z0-9][a-z0-9.]*[a-z0-9]@[a-z0-9.-]+\.[a-z]{2,}$`,
)

func (s *AuthenticationService) RegisterUser(
	ctx context.Context,
	email string,
	password string,
	fullName string,
) (domain.User, error) {
	email = normalizeEmail(email)

	if err := validateEmail(email); err != nil {
      return domain.User{}, err
  	}
	
	if err := validateFullName(fullName); err != nil {
		return domain.User{}, err
	}

	if err := validatePassword(password); err != nil {
		return domain.User{}, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("hash password: %w", err)
	}

	createdUser, err := s.authenticationRepository.RegisterUser(
		ctx,
		email,
		string(passwordHash),
		fullName,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf(
			"create user with email %q: %w",
			email,
			err,
		)
	}

	return createdUser, nil
}

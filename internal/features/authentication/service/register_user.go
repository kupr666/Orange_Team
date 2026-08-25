package auth_service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
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

func (s *AuthService) RegisterUser(
	ctx context.Context,
	email string,
	password string,
	fullName string,
) (domain.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
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

	createdUser, err := s.authRepository.RegisterUser(
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

func validateEmail(email string) error {
	emailLength := utf8.RuneCountInString(email)
	if emailLength < minimumEmailLength || emailLength > maximumEmailLength {
		return fmt.Errorf(
			"email length must be between %d and %d characters: %w",
			minimumEmailLength,
			maximumEmailLength,
			core_errors.ErrInvalidArgument,
		)
	}

	if strings.Contains(email, "..") || !emailPattern.MatchString(email) {
		return fmt.Errorf("invalid email format: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

func validateFullName(fullName string) error {
	fullNameLength := utf8.RuneCountInString(fullName)
	if fullName != strings.TrimSpace(fullName) ||
		fullNameLength < minimumFullNameLength ||
		fullNameLength > maximumFullNameLength {
		return fmt.Errorf(
			"full name must not have surrounding spaces and must contain between %d and %d characters: %w",
			minimumFullNameLength,
			maximumFullNameLength,
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

func validatePassword(password string) error {
	passwordLength := len([]byte(password))
	if passwordLength < minimumPasswordLengthBytes ||
		passwordLength > maximumPasswordLengthBytes {
		return fmt.Errorf(
			"password length must be between %d and %d bytes: %w",
			minimumPasswordLengthBytes,
			maximumPasswordLengthBytes,
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

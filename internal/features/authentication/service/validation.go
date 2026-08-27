package authentication_service

import (
	"fmt"
	"strings"
	"unicode/utf8"

	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
)

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
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

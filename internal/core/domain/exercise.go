package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
)

var exerciseTextPattern = regexp.MustCompile(`^[A-Za-zА-Яа-яёЁ ]+$`)

type Exercise struct {
	ID          int
	Version     int64
	Name        string
	Description string
	Difficulty  int
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	Type        string
}

func NewExercise(
	id int,
	version int64,
	name string,
	description string,
	difficulty int,
	createdAt time.Time,
	updatedAt *time.Time,
	exerciseType string,
) Exercise {
	return Exercise{
		ID:          id,
		Version:     version,
		Name:        name,
		Description: description,
		Difficulty:  difficulty,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Type:        exerciseType,
	}
}

func (e Exercise) Validate() error {
	if err := validateExerciseText("name", e.Name, 3, 100); err != nil {
		return err
	}

	if err := validateExerciseText("description", e.Description, 1, 1000); err != nil {
		return err
	}

	if e.Difficulty < 1 || e.Difficulty > 10 {
		return fmt.Errorf(
			"exercise difficulty must be between 1 and 10: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if e.Type != "weight" && e.Type != "duration" {
		return fmt.Errorf(
			"exercise type must be either weight or duration: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

func validateExerciseText(field, value string, minLength, maxLength int) error {
	if value != strings.TrimSpace(value) {
		return fmt.Errorf(
			"exercise %s must not have leading or trailing spaces: %w",
			field,
			core_errors.ErrInvalidArgument,
		)
	}

	length := utf8.RuneCountInString(value)
	if length < minLength || length > maxLength {
		return fmt.Errorf(
			"exercise %s length must be between %d and %d characters: %w",
			field,
			minLength,
			maxLength,
			core_errors.ErrInvalidArgument,
		)
	}

	if !exerciseTextPattern.MatchString(value) {
		return fmt.Errorf(
			"exercise %s must contain only Latin or Cyrillic letters and spaces: %w",
			field,
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

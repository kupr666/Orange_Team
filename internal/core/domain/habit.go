package domain

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
)

const (
	HabitNameMinLength        = 2
	HabitNameMaxLength        = 80
	HabitDescriptionMaxLength = 500
)

type Habit struct {
	ID                uuid.UUID
	Version           int64
	UserID            uuid.UUID
	Name              string
	Description       string
	CurrentStreak     int
	LastCompletedDate *time.Time
	CompletedToday    bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func NewHabit(
	id uuid.UUID,
	version int64,
	userID uuid.UUID,
	name string,
	description string,
	currentStreak int,
	lastCompletedDate *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) Habit {
	return Habit{
		ID:                id,
		Version:           version,
		UserID:            userID,
		Name:              name,
		Description:       description,
		CurrentStreak:     currentStreak,
		LastCompletedDate: lastCompletedDate,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}
}

func (h Habit) ValidateForCreation() error {
	if h.UserID == uuid.Nil {
		return fmt.Errorf("habit user ID is empty: %w", core_errors.ErrInvalidArgument)
	}

	if h.Name != strings.TrimSpace(h.Name) {
		return fmt.Errorf(
			"habit name must not have leading or trailing spaces: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	nameLength := utf8.RuneCountInString(h.Name)
	if nameLength < HabitNameMinLength || nameLength > HabitNameMaxLength {
		return fmt.Errorf(
			"habit name length must be between %d and %d characters: %w",
			HabitNameMinLength,
			HabitNameMaxLength,
			core_errors.ErrInvalidArgument,
		)
	}

	if utf8.RuneCountInString(h.Description) > HabitDescriptionMaxLength {
		return fmt.Errorf(
			"habit description must not exceed %d characters: %w",
			HabitDescriptionMaxLength,
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

func (h Habit) Complete(on time.Time) (Habit, bool) {
	completedOn := startOfUTCDay(on)
	if h.LastCompletedDate != nil && sameUTCDay(*h.LastCompletedDate, completedOn) {
		return h.ViewAt(completedOn), false
	}

	if h.LastCompletedDate != nil && sameUTCDay(*h.LastCompletedDate, completedOn.AddDate(0, 0, -1)) {
		h.CurrentStreak++
	} else {
		h.CurrentStreak = 1
	}

	h.LastCompletedDate = &completedOn
	h.CompletedToday = true
	return h, true
}

func (h Habit) ViewAt(on time.Time) Habit {
	today := startOfUTCDay(on)
	if h.LastCompletedDate == nil {
		h.CurrentStreak = 0
		h.CompletedToday = false
		return h
	}

	h.CompletedToday = sameUTCDay(*h.LastCompletedDate, today)
	yesterday := today.AddDate(0, 0, -1)
	if !h.CompletedToday && !sameUTCDay(*h.LastCompletedDate, yesterday) {
		h.CurrentStreak = 0
	}

	return h
}

func startOfUTCDay(value time.Time) time.Time {
	utc := value.UTC()
	return time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
}

func sameUTCDay(first time.Time, second time.Time) bool {
	return startOfUTCDay(first).Equal(startOfUTCDay(second))
}

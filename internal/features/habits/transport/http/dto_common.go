package habits_transport_http

import (
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

type HabitDTOResponse struct {
	ID                uuid.UUID `json:"id"`
	Version           int64     `json:"version"`
	UserID            uuid.UUID `json:"user_id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	CurrentStreak     int       `json:"current_streak"`
	LastCompletedDate *string   `json:"last_completed_date"`
	CompletedToday    bool      `json:"completed_today"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func habitDTOFromDomain(habit domain.Habit) HabitDTOResponse {
	var lastCompletedDate *string
	if habit.LastCompletedDate != nil {
		formatted := habit.LastCompletedDate.UTC().Format("2006-01-02")
		lastCompletedDate = &formatted
	}

	return HabitDTOResponse{
		ID:                habit.ID,
		Version:           habit.Version,
		UserID:            habit.UserID,
		Name:              habit.Name,
		Description:       habit.Description,
		CurrentStreak:     habit.CurrentStreak,
		LastCompletedDate: lastCompletedDate,
		CompletedToday:    habit.CompletedToday,
		CreatedAt:         habit.CreatedAt,
		UpdatedAt:         habit.UpdatedAt,
	}
}

func habitsDTOFromDomain(habits []domain.Habit) []HabitDTOResponse {
	response := make([]HabitDTOResponse, len(habits))
	for i, habit := range habits {
		response[i] = habitDTOFromDomain(habit)
	}
	return response
}

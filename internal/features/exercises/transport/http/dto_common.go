package exercises_transport_http

import (
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

type ExerciseDTOResponse struct {
	ID          uuid.UUID  `json:"id"`
	Version     int        `json:"version"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Difficulty  int        `json:"difficulty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	Type        string     `json:"type"`
}

func exerciseDTOFromDomain(exercise domain.Exercise) ExerciseDTOResponse {
	return ExerciseDTOResponse{
		ID:          exercise.ID,
		Version:     exercise.Version,
		Name:        exercise.Name,
		Description: exercise.Description,
		Difficulty:  exercise.Difficulty,
		CreatedAt:   exercise.CreatedAt,
		UpdatedAt:   exercise.UpdatedAt,
		Type:        exercise.Type,
	}
}

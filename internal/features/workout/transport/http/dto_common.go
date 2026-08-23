package workouts_transport_http

import (
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

type WorkoutResponseDTO struct {
	ID                       uuid.UUID  `json:"id"`
	Version                  int64      `json:"version"`
	UserID                   uuid.UUID  `json:"user_id"`
	Status                   string     `json:"status"`
	StartedAt                *time.Time `json:"started_at"`
	CompletedAt              *time.Time `json:"completed_at"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
	WorkoutScore             int        `json:"workout_score"`
	Intensity                *int       `json:"intensity"`
	PersonalScoreCoefficient int        `json:"personal_score_coefficient"`
}

func workoutDTOFromDomain(
	workout domain.Workout,
) WorkoutResponseDTO {
	return WorkoutResponseDTO{
		ID:                       workout.ID,
		Version:                  workout.Version,
		UserID:                   workout.UserID,
		Status:                   workout.Status,
		StartedAt:                workout.StartedAt,
		CompletedAt:              workout.CompletedAt,
		CreatedAt:                workout.CreatedAt,
		UpdatedAt:                workout.UpdatedAt,
		WorkoutScore:             workout.WorkoutScore,
		Intensity:                workout.Intensity,
		PersonalScoreCoefficient: workout.PersonalScoreCoefficient,
	}
}

func workoutDTOsFromDomains(
	workouts []domain.Workout,
) []WorkoutResponseDTO {
	response := make([]WorkoutResponseDTO, len(workouts))

	for i, workout := range workouts {
		response[i] = workoutDTOFromDomain(workout)
	}

	return response
}
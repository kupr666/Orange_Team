package domain

import (
	"time"

	"github.com/google/uuid"
)

type Workout struct {
	ID                       uuid.UUID
	Version                  int
	UserID                   uuid.UUID
	Status                   string
	StartedAt                *time.Time
	CompletedAt              *time.Time
	CreatedAt                time.Time
	UpdatedAt                time.Time
	WorkoutScore             int
	Intensity                *int
	PersonalScoreCoefficient int
}

func NewWorkout(
	id uuid.UUID,
	version int,
	userID uuid.UUID,
	status string,
	startedAt *time.Time,
	completedAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
	workoutScore int,
	intensity *int,
	personalScoreCoefficient int,
) Workout {
	return Workout{
		ID:                       id,
		Version:                  version,
		UserID:                   userID,
		Status:                   status,
		StartedAt:                startedAt,
		CompletedAt:              completedAt,
		CreatedAt:                createdAt,
		UpdatedAt:                updatedAt,
		WorkoutScore:             workoutScore,
		Intensity:                intensity,
		PersonalScoreCoefficient: personalScoreCoefficient,
	}
}

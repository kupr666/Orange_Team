package domain

import (
	"time"

	"github.com/google/uuid"
)

type Workout struct {
	ID uuid.UUID
	Version int64
	UserID int // UUID ?
	Status string
	Started_at *time.Time
	Completed_at *time.Time
	Created_at time.Time
	Updated_at *time.Time
	WorkoutScore int
	Intensity int
	PersonalScoreCoefficient int
}

func NewWorkout(
	id uuid.UUID,
	version int64,
	userID int,
	status string,
)
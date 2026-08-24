package domain

import (
	"time"

	"github.com/google/uuid"
)

type CreatedExercise struct {
	ID           uuid.UUID
	Version      int
	ExerciseID   uuid.UUID
	WorkoutID    uuid.UUID
	Weight       *int
	Sets         *int
	Reps         *int
	Duration     *int
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	Completed    bool
	ExerciseLoad int
}

func NewCreatedExercise(
	id uuid.UUID,
	version int,
	exerciseID uuid.UUID,
	workoutID uuid.UUID,
	weight *int,
	sets *int,
	reps *int,
	duration *int,
	createdAt time.Time,
	updatedAt *time.Time,
	completed bool,
	exerciseLoad int,
) CreatedExercise {
	return CreatedExercise{
		ID:           id,
		Version:      version,
		ExerciseID:   exerciseID,
		WorkoutID:    workoutID,
		Weight:       weight,
		Sets:         sets,
		Reps:         reps,
		Duration:     duration,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		Completed:    completed,
		ExerciseLoad: exerciseLoad,
	}
}

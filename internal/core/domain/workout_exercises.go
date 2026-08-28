package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
)

const (
	ExerciseTypeWeight   = "weight"
	ExerciseTypeDuration = "duration"
)

type WorkoutExercise struct {
	ID           uuid.UUID
	Version      int
	WorkoutID    uuid.UUID
	ExerciseID   uuid.UUID
	Weight       *int
	Sets         *int
	Reps         *int
	Duration     *int
	Completed    bool
	ExerciseLoad int
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

func NewWorkoutExercise(
	id uuid.UUID,
	Version int,
	workoutID uuid.UUID,
	exerciseID uuid.UUID,
	weight *int,
	sets *int,
	reps *int,
	duration *int,
	completed bool,
	exerciseLoad int,
	createdAt time.Time,
	updatedAt *time.Time,
) WorkoutExercise {
	return WorkoutExercise{
		ID:           id,
		Version:      Version,
		WorkoutID:    workoutID,
		ExerciseID:   exerciseID,
		Weight:       weight,
		Sets:         sets,
		Reps:         reps,
		Duration:     duration,
		Completed:    completed,
		ExerciseLoad: exerciseLoad,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}

func CreateWorkoutExercise(
	workoutID uuid.UUID,
	exerciseID uuid.UUID,
	weight *int,
	sets *int,
	reps *int,
	duration *int,
	completed bool,
	exerciseType string,
) WorkoutExercise {
	var (
		id           = uuid.New()
		version      = 1
		createdAt    = time.Now()
		updatedAt    = &createdAt
		exerciseLoad = 0
	)

	workoutExercise := NewWorkoutExercise(
		id,
		version,
		workoutID,
		exerciseID,
		weight,
		sets,
		reps,
		duration,
		completed,
		exerciseLoad,
		createdAt,
		updatedAt,
	)

	workoutExercise.ExerciseLoad = workoutExercise.CalculateLoad(exerciseType)
	return workoutExercise
}

func (workoutExercise *WorkoutExercise) Validate() error {
	hasAnyWeightField := workoutExercise.Weight != nil || workoutExercise.Sets != nil || workoutExercise.Reps != nil
	hasAllWeightFields := workoutExercise.Weight != nil && workoutExercise.Sets != nil && workoutExercise.Reps != nil
	hasDuration := workoutExercise.Duration != nil

	if hasAnyWeightField && !hasAllWeightFields {
		return fmt.Errorf(
			"weight, sets and reps must be provided together: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	if !hasAllWeightFields && !hasDuration {
		return fmt.Errorf(
			"either (weight, sets, reps) or duration must be provided: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	if hasAllWeightFields && hasDuration {
		return fmt.Errorf(
			"cannot provide both (weight, sets, reps) and duration: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	if workoutExercise.Weight != nil && *workoutExercise.Weight < 0 {
		return fmt.Errorf(
			"weight must be >= 0: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	if workoutExercise.Sets != nil && *workoutExercise.Sets < 0 {
		return fmt.Errorf(
			"sets must be >= 0: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	if workoutExercise.Reps != nil && *workoutExercise.Reps < 0 {
		return fmt.Errorf(
			"reps must be >= 0: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	if workoutExercise.Duration != nil && *workoutExercise.Duration < 0 {
		return fmt.Errorf(
			"duration must be >= 0: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	return nil
}

func (workoutExercise *WorkoutExercise) CalculateLoad(exerciseType string) int {
	if exerciseType == ExerciseTypeWeight && workoutExercise.Weight != nil && workoutExercise.Sets != nil && workoutExercise.Reps != nil {
		return *workoutExercise.Weight * *workoutExercise.Sets * *workoutExercise.Reps
	}
	if exerciseType == ExerciseTypeDuration && workoutExercise.Duration != nil {
		return *workoutExercise.Duration
	}
	return 0
}

type WorkoutExercisePatch struct {
	Weight    Nullable[int]
	Sets      Nullable[int]
	Reps      Nullable[int]
	Duration  Nullable[int]
	Completed Nullable[bool]
}

func NewWorkoutExercisePatch(
	weight Nullable[int],
	sets Nullable[int],
	reps Nullable[int],
	duration Nullable[int],
	completed Nullable[bool],
) WorkoutExercisePatch {
	return WorkoutExercisePatch{
		Weight:    weight,
		Sets:      sets,
		Reps:      reps,
		Duration:  duration,
		Completed: completed,
	}
}

func (p *WorkoutExercisePatch) Validate() error {
	if !p.Weight.Set && !p.Sets.Set && !p.Reps.Set && !p.Duration.Set && !p.Completed.Set {
		return fmt.Errorf(
			"at least one field must be patched: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.Completed.Set && p.Completed.Value == nil {
		return fmt.Errorf(
			"completed cannot be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	return nil
}

func (workoutExercise *WorkoutExercise) ApplyPatch(patch WorkoutExercisePatch, exerciseType string) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf(
			"validate workout exercise patch: %w",
			err,
		)
	}

	tmp := *workoutExercise

	if patch.Weight.Set {
		tmp.Weight = patch.Weight.Value
	}
	if patch.Sets.Set {
		tmp.Sets = patch.Sets.Value
	}
	if patch.Reps.Set {
		tmp.Reps = patch.Reps.Value
	}
	if patch.Duration.Set {
		tmp.Duration = patch.Duration.Value
	}
	if patch.Completed.Set {
		tmp.Completed = *patch.Completed.Value
	}

	tmp.ExerciseLoad = tmp.CalculateLoad(exerciseType)

	if err := tmp.ValidateForExerciseType(exerciseType); err != nil {
		return fmt.Errorf(
			"validate patched workout exercise: %w",
			err,
		)
	}

	*workoutExercise = tmp
	return nil
}

func (workoutExercise *WorkoutExercise) ValidateForExerciseType(exerciseType string) error {
	if err := workoutExercise.Validate(); err != nil {
		return err
	}

	switch exerciseType {
	case ExerciseTypeWeight:
		if workoutExercise.Weight == nil ||
			workoutExercise.Sets == nil ||
			workoutExercise.Reps == nil ||
			workoutExercise.Duration != nil {
			return fmt.Errorf(
				"weight exercise requires weight, sets and reps only: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	case ExerciseTypeDuration:
		if workoutExercise.Duration == nil ||
			workoutExercise.Weight != nil ||
			workoutExercise.Sets != nil ||
			workoutExercise.Reps != nil {
			return fmt.Errorf(
				"duration exercise requires duration only: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	default:
		return fmt.Errorf(
			"unsupported exercise type %q: %w",
			exerciseType,
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

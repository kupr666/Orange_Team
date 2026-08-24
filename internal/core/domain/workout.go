package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
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

func (w *Workout) Validate() error {
	allowed := map[string]bool{
		"planned": true, "in_progress": true,
		"completed": true, "cancelled": true,
	}
	if !allowed[w.Status] {
		return fmt.Errorf("invalid status: %s: %w", w.Status, core_errors.ErrInvalidArgument)
	}

	switch w.Status {
	case "planned":
		if w.StartedAt != nil || w.CompletedAt != nil {
			return fmt.Errorf("planned workout must have started_at and completed_at NULL: %w", core_errors.ErrInvalidArgument)
		}
	case "in_progress":
		if w.StartedAt == nil || w.CompletedAt != nil {
			return fmt.Errorf("in_progress workout must have started_at not NULL and completed_at NULL: %w", core_errors.ErrInvalidArgument)
		}
	case "completed":
		if w.StartedAt == nil || w.CompletedAt == nil {
			return fmt.Errorf("completed workout must have started_at and completed_at not NULL: %w", core_errors.ErrInvalidArgument)
		}
		if w.CompletedAt.Before(*w.StartedAt) {
			return fmt.Errorf("completed_at must be after started_at: %w", core_errors.ErrInvalidArgument)
		}
		if w.Intensity == nil {
			return fmt.Errorf("completed workout must have intensity set: %w", core_errors.ErrInvalidArgument)
		}
		if *w.Intensity < 1 || *w.Intensity > 10 {
			return fmt.Errorf("intensity must be between 1 and 10: %w", core_errors.ErrInvalidArgument)
		}
	case "cancelled":
		if w.CompletedAt != nil {
			return fmt.Errorf("cancelled workout must have completed_at NULL: %w", core_errors.ErrInvalidArgument)
		}
	}

	if w.Intensity != nil && (*w.Intensity < 1 || *w.Intensity > 10) {
		return fmt.Errorf("intensity must be between 1 and 10: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

type WorkoutPatch struct {
	Status      Nullable[string]
	StartedAt   Nullable[time.Time]
	CompletedAt Nullable[time.Time]
	Intensity   Nullable[int]
}

func NewWorkoutPatch(
	status Nullable[string],
	startedAt Nullable[time.Time],
	completedAt Nullable[time.Time],
	intensity Nullable[int],
) WorkoutPatch {
	return WorkoutPatch{
		Status:      status,
		StartedAt:   startedAt,
		CompletedAt: completedAt,
		Intensity:   intensity,
	}
}

func (p *WorkoutPatch) Validate() error {
	if p.Status.Set && p.Status.Value == nil {
		return fmt.Errorf(
			"status cannot be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	return nil
}

func (w *Workout) ApplyPatch(patch WorkoutPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf(
			"validate workout patch: %w",
			err,
		)
	}

	tmp := *w
	newStatus := tmp.Status
	if patch.Status.Set {
		newStatus = *patch.Status.Value
	}

	if tmp.Status == "cancelled" {
		if patch.Status.Set || patch.StartedAt.Set || patch.CompletedAt.Set || patch.Intensity.Set {
			return fmt.Errorf(
				"cancelled workout cannot be modified: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}
	if tmp.Status == "completed" {
		if patch.Status.Set || patch.StartedAt.Set || patch.CompletedAt.Set {
			return fmt.Errorf(
				"completed workout cannot change status or dates: %w",
				core_errors.ErrInvalidArgument,
			)
		}
		if patch.Intensity.Set && patch.Intensity.Value == nil {
			return fmt.Errorf(
				"intensity cannot be set to NULL for completed workout: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if patch.Status.Set {
		if tmp.Status == "completed" || tmp.Status == "cancelled" {
			return fmt.Errorf(
				"cannot change status from terminal status %s: %w",
				tmp.Status,
				core_errors.ErrInvalidArgument,
			)
		}
		switch tmp.Status {
		case "planned":
			if newStatus != "in_progress" && newStatus != "cancelled" {
				return fmt.Errorf(
					"from planned can only go to in_progress or cancelled: %w",
					core_errors.ErrInvalidArgument,
				)
			}
		case "in_progress":
			if newStatus != "completed" {
				return fmt.Errorf(
					"from in_progress can only go to completed: %w",
					core_errors.ErrInvalidArgument,
				)
			}
		}
	}

	if patch.Intensity.Set {
		if newStatus != "completed" && tmp.Status != "completed" {
			return fmt.Errorf(
				"intensity can only be set for completed workouts: %w",
				core_errors.ErrInvalidArgument,
			)
		}
		if patch.Status.Set && newStatus == "completed" && patch.Intensity.Value == nil {
			return fmt.Errorf(
				"intensity is required when moving to completed: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	} else {
		if patch.Status.Set && newStatus == "completed" {
			return fmt.Errorf(
				"intensity is required when moving to completed: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if patch.Status.Set {
		tmp.Status = newStatus
	}
	if patch.StartedAt.Set {
		tmp.StartedAt = patch.StartedAt.Value
	}
	if patch.CompletedAt.Set {
		tmp.CompletedAt = patch.CompletedAt.Value
	}
	if patch.Intensity.Set {
		tmp.Intensity = patch.Intensity.Value
	}

	if patch.Status.Set {
		switch tmp.Status {
		case "in_progress":
			if tmp.StartedAt == nil {
				now := time.Now()
				tmp.StartedAt = &now
			}
		case "completed":
			if tmp.CompletedAt == nil {
				now := time.Now()
				tmp.CompletedAt = &now
			}
		}
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf(
			"validate patched workout: %w",
			err,
		)
	}

	*w = tmp
	return nil
}

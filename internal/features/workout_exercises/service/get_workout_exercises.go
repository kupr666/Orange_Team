package workout_exercises_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
)

func (s *WorkoutExercisesService) GetWorkoutExercises(
	ctx context.Context,
	userID uuid.UUID,
	workoutID uuid.UUID,
) ([]domain.WorkoutExercise, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf(
			"user ID is empty: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	workout, err := s.repository.GetWorkout(ctx, userID, workoutID)
	if err != nil {
		return nil, fmt.Errorf(
			"get workout: %w",
			err,
		)
	}
	if workout.UserID != userID {
		return nil, fmt.Errorf(
			"user is not the owner of the workout: %w",
			core_errors.ErrForbidden,
		)
	}

	workoutExercises, err := s.repository.GetWorkoutExercises(ctx, workoutID)
	if err != nil {
		return nil, fmt.Errorf(
			"get workout exercises from repository: %w",
			err,
		)
	}

	return workoutExercises, nil
}

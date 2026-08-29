package workout_exercises_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
)

func (s *WorkoutExercisesService) DeleteWorkoutExercise(
	ctx context.Context,
	userID uuid.UUID,
	workoutID uuid.UUID,
	workoutExerciseID uuid.UUID,
) error {
	if userID == uuid.Nil {
		return fmt.Errorf(
			"user ID is empty: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	workout, err := s.workoutRepository.GetWorkout(ctx, userID, workoutID)
	if err != nil {
		return fmt.Errorf(
			"get workout: %w",
			err,
		)
	}

	if workout.UserID != userID {
		return fmt.Errorf(
			"user is not the owner of the workout: %w",
			core_errors.ErrForbidden,
		)
	}

	if !workout.CanModifyWorkoutExercise() {
		return fmt.Errorf(
			"can't delete exercise when workout status is %s: %w",
			workout.Status,
			core_errors.ErrConflict,
		)
	}

	if err := s.workoutExercisesRepository.DeleteWorkoutExercise(ctx, workoutID, workoutExerciseID); err != nil {
		return fmt.Errorf(
			"delete workout exercise: %w",
			err,
		)
	}

	if err := s.recalculateScore(ctx, workoutID, userID); err != nil {
		return fmt.Errorf(
			"recalculate workout score: %w",
			err,
		)
	}

	return nil
}

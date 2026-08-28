// delete_workout_exercise.go (сервис)
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

	workout, err := s.workoutRepository.GetWorkout(ctx, workoutID)
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

	if err := s.workoutExercisesRepository.DeleteWorkoutExercise(ctx, workoutID, workoutExerciseID); err != nil {
		return fmt.Errorf(
			"delete workout exercise: %w",
			err,
		)
	}

	// Пересчёт workout_score (пока заглушка)
	_ = s.recalculateScore(ctx, workoutID)

	return nil
}

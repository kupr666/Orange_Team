package workout_exercises_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
)

func (s *WorkoutExercisesService) PatchWorkoutExercise(
	ctx context.Context,
	userID uuid.UUID,
	workoutID uuid.UUID,
	workoutExerciseID uuid.UUID,
	patch domain.WorkoutExercisePatch,
) (domain.WorkoutExercise, error) {
	if userID == uuid.Nil {
		return domain.WorkoutExercise{}, fmt.Errorf(
			"user ID is empty: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	workout, err := s.repository.GetWorkout(ctx, userID, workoutID)
	if err != nil {
		return domain.WorkoutExercise{}, fmt.Errorf(
			"get workout: %w",
			err,
		)
	}
	if workout.UserID != userID {
		return domain.WorkoutExercise{}, fmt.Errorf(
			"user is not the owner of the workout: %w",
			core_errors.ErrForbidden,
		)
	}

	workoutExercise, err := s.repository.GetWorkoutExercise(
		ctx,
		workoutID,
		workoutExerciseID,
	)
	if err != nil {
		return domain.WorkoutExercise{}, fmt.Errorf(
			"get workout exercise: %w",
			err,
		)
	}

	exercise, err := s.repository.GetExercise(ctx, workoutExercise.ExerciseID)
	if err != nil {
		return domain.WorkoutExercise{}, fmt.Errorf(
			"get exercise: %w",
			err,
		)
	}

	if err := workoutExercise.ApplyPatch(patch, exercise.Type); err != nil {
		return domain.WorkoutExercise{}, fmt.Errorf(
			"apply workout exercise patch: %w",
			err,
		)
	}

	updatedWorkoutExercise, err := s.repository.PatchWorkoutExercise(ctx, workoutExercise)
	if err != nil {
		return domain.WorkoutExercise{}, fmt.Errorf(
			"patch workout exercise: %w",
			err,
		)
	}

	if err := s.recalculateScore(ctx, workoutID); err != nil {
		return domain.WorkoutExercise{}, fmt.Errorf(
			"recalculate workout score: %w",
			err,
		)
	}

	return updatedWorkoutExercise, nil
}

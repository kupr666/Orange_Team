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

	workout, err := s.workoutRepository.GetWorkout(ctx, userID, workoutID)
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

	if !workout.CanModifyWorkoutExercise() {
		return domain.WorkoutExercise{}, fmt.Errorf(
			"can't patch exercise when workout status is %s: %w",
			workout.Status,
			core_errors.ErrConflict,
		)
	}

	if patch.Completed.Set && patch.Completed.Value != nil && *patch.Completed.Value {
		if workout.Status != domain.StatusInProgress {
			return domain.WorkoutExercise{}, fmt.Errorf(
				"can only mark exercise as completed when workout is in_progress: %w",
				core_errors.ErrConflict,
			)
		}
	}

	currentWorkoutExercise, err := s.workoutExercisesRepository.GetWorkoutExercise(ctx, workoutID, workoutExerciseID)
	if err != nil {
		return domain.WorkoutExercise{}, fmt.Errorf(
			"get workout exercise: %w",
			err,
		)
	}

	exercise, err := s.exerciseRepository.GetExercise(ctx, currentWorkoutExercise.ExerciseID)
	if err != nil {
		return domain.WorkoutExercise{}, fmt.Errorf(
			"get exercise type: %w",
			err,
		)
	}

	if err := currentWorkoutExercise.ApplyPatch(patch, exercise.Type); err != nil {
		return domain.WorkoutExercise{}, fmt.Errorf(
			"apply patch: %w",
			err,
		)
	}

	if err := currentWorkoutExercise.ValidateForWorkoutExerciseType(exercise.Type); err != nil {
		return domain.WorkoutExercise{}, fmt.Errorf(
			"validate patched exercise: %w",
			err,
		)
	}

	patchedWorkoutExercise, err := s.workoutExercisesRepository.PatchWorkoutExercise(ctx, currentWorkoutExercise)
	if err != nil {
		return domain.WorkoutExercise{}, fmt.Errorf(
			"update workout exercise: %w",
			err,
		)
	}

	if err := s.recalculateScore(ctx, workoutID, userID); err != nil {
		return domain.WorkoutExercise{}, fmt.Errorf(
			"recalculate workout score: %w",
			err,
		)
	}

	return patchedWorkoutExercise, nil
}

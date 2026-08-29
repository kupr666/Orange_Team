package workouts_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
)

func (s *WorkoutsService) PatchWorkout(
	ctx context.Context,
	userID uuid.UUID,
	workoutID uuid.UUID,
	patch domain.WorkoutPatch,
) (domain.Workout, error) {
	workout, err := s.workoutsRepository.GetWorkout(ctx, userID, workoutID)
	if err != nil {
		return domain.Workout{}, fmt.Errorf(
			"get workout from repository: %w",
			err,
		)
	}

	if err := workout.ApplyPatch(patch); err != nil {
		return domain.Workout{}, fmt.Errorf(
			"apply workout patch: %w",
			err,
		)
	}

	if workout.Status == domain.StatusCompleted {
		workoutExercises, err := s.workoutExerciseReader.GetWorkoutExercises(ctx, workoutID)
		if err != nil {
			return domain.Workout{}, fmt.Errorf(
				"get exercises for workout: %w",
				err,
			)
		}

		completedExercises := 0
		for _, workoutExercise := range workoutExercises {
			if workoutExercise.Completed {
				completedExercises++
			}
		}
		if completedExercises < 1 {
			return domain.Workout{}, fmt.Errorf(
				"cannot complete workout: at least one exercise must be completed: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	updatedWorkout, err := s.workoutsRepository.PatchWorkout(ctx, userID, workout)
	if err != nil {
		return domain.Workout{}, fmt.Errorf(
			"update workout in repository: %w",
			err,
		)
	}

	if updatedWorkout.Status == domain.StatusCompleted {
		if err := s.userRepository.UpdateUserWorkoutScore(ctx, userID); err != nil {
			return domain.Workout{}, fmt.Errorf(
				"update user workout score: %w",
				err,
			)
		}
	}

	return updatedWorkout, nil
}

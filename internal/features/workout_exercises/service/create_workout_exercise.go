package workout_exercises_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
)

func (s *WorkoutExercisesService) CreateWorkoutExercise(
	ctx context.Context,
	userID uuid.UUID,
	workoutID uuid.UUID,
	exerciseID uuid.UUID,
	weight *int,
	sets *int,
	reps *int,
	duration *int,
	completed bool,
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
			"can't create exercise when workout status is %s: %w",
			workout.Status,
			core_errors.ErrConflict,
		)
	}

	if completed && workout.Status != domain.StatusInProgress {
		return domain.WorkoutExercise{}, fmt.Errorf(
			"can only create exercise as completed when workout is in_progress: %w",
			core_errors.ErrConflict,
		)
	}

	exercise, err := s.exerciseRepository.GetExercise(ctx, exerciseID)
	if err != nil {
		return domain.WorkoutExercise{}, fmt.Errorf(
			"get exercise: %w",
			err,
		)
	}

	if exercise.Type == domain.ExerciseTypeWeight && (weight == nil || sets == nil || reps == nil) {
		return domain.WorkoutExercise{}, fmt.Errorf(
			"weight exercise requires weight, sets, reps: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if exercise.Type == domain.ExerciseTypeDuration && duration == nil {
		return domain.WorkoutExercise{}, fmt.Errorf(
			"duration exercise requires duration: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	workoutExercise := domain.CreateWorkoutExercise(
		workoutID,
		exerciseID,
		weight,
		sets,
		reps,
		duration,
		completed,
		exercise.Type,
	)

	if err := workoutExercise.ValidateForWorkoutExerciseType(exercise.Type); err != nil {
		return domain.WorkoutExercise{}, fmt.Errorf(
			"validate workout exercise: %w",
			err,
		)
	}

	workoutExercise, err = s.workoutExercisesRepository.CreateWorkoutExercise(ctx, workoutExercise)
	if err != nil {
		return domain.WorkoutExercise{}, fmt.Errorf(
			"create workout exercise: %w",
			err,
		)
	}

	if err := s.recalculateScore(ctx, workoutID, userID); err != nil {
		return domain.WorkoutExercise{}, fmt.Errorf(
			"recalculate workout score: %w",
			err,
		)
	}

	return workoutExercise, nil
}

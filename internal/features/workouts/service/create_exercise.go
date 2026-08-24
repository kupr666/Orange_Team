package workouts_service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
)

func (s *WorkoutsService) CreateExercise(
	ctx context.Context,
	userID uuid.UUID,
	workoutID uuid.UUID,
	exerciseID uuid.UUID,
	weight *int,
	sets *int,
	reps *int,
	duration *int,
) (domain.CreatedExercise, error) {
	if err := validateCreateExercise(
		userID,
		workoutID,
		exerciseID,
		weight,
		sets,
		reps,
		duration,
	); err != nil {
		return domain.CreatedExercise{}, err
	}

	// TODO: calculate exercise load when the scoring formula is defined.
	exerciseLoad := 0

	exercise := domain.NewCreatedExercise(
		uuid.New(),
		1,
		exerciseID,
		workoutID,
		weight,
		sets,
		reps,
		duration,
		time.Now().UTC(),
		nil,
		false,
		exerciseLoad,
	)

	createdExercise, err := s.workoutsRepository.CreateExercise(
		ctx,
		userID,
		exercise,
	)
	if err != nil {
		return domain.CreatedExercise{}, fmt.Errorf(
			"create exercise in repository: %w",
			err,
		)
	}

	return createdExercise, nil
}

func validateCreateExercise(
	userID uuid.UUID,
	workoutID uuid.UUID,
	exerciseID uuid.UUID,
	weight *int,
	sets *int,
	reps *int,
	duration *int,
) error {
	if err := validateCreateExerciseIDs(userID, workoutID, exerciseID); err != nil {
		return err
	}

	if err := validateExerciseParameters(weight, sets, reps, duration); err != nil {
		return err
	}

	return validateExerciseParameterValues(weight, sets, reps, duration)
}

func validateCreateExerciseIDs(
	userID uuid.UUID,
	workoutID uuid.UUID,
	exerciseID uuid.UUID,
) error {
	if userID == uuid.Nil {
		return fmt.Errorf(
			"user ID is empty: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if workoutID == uuid.Nil {
		return fmt.Errorf(
			"workout ID is empty: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if exerciseID == uuid.Nil {
		return fmt.Errorf(
			"exercise ID is empty: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

func validateExerciseParameters(
	weight *int,
	sets *int,
	reps *int,
	duration *int,
) error {
	isStrengthExercise := weight != nil &&
		sets != nil &&
		reps != nil &&
		duration == nil

	isDurationExercise := weight == nil &&
		sets == nil &&
		reps == nil &&
		duration != nil

	if !isStrengthExercise && !isDurationExercise {
		return fmt.Errorf(
			"provide either weight, sets and reps or duration: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

func validateExerciseParameterValues(
	weight *int,
	sets *int,
	reps *int,
	duration *int,
) error {

	if weight != nil && *weight < 0 {
		return fmt.Errorf(
			"weight must not be negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if sets != nil && *sets <= 0 {
		return fmt.Errorf(
			"sets must be positive: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if reps != nil && *reps <= 0 {
		return fmt.Errorf(
			"reps must be positive: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if duration != nil && *duration <= 0 {
		return fmt.Errorf(
			"duration must be positive: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

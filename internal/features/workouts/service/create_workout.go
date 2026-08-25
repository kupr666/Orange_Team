package workouts_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
)

func (s *WorkoutsService) CreateWorkout(
	ctx context.Context,
	userID uuid.UUID,
) (domain.Workout, error) {
	if err := validateCreateWorkout(userID); err != nil {
		return domain.Workout{}, err
	}

	workout, err := s.workoutsRepository.CreateWorkout(ctx, userID)

	if err != nil {
		return domain.Workout{}, fmt.Errorf(
			"create workout in repository: %w",
			err,
		)
	}

	return workout, nil
}

func validateCreateWorkout(userID uuid.UUID) error {
	if userID == uuid.Nil {
		return fmt.Errorf(
			"user ID is empty: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

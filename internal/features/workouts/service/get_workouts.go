package workouts_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

func (s *WorkoutsService) GetWorkouts(
	ctx context.Context,
	userID uuid.UUID,
) ([]domain.Workout, error) {

	workouts, err := s.workoutsRepository.GetWorkouts(ctx, userID)
	if err != nil {
		return []domain.Workout{}, fmt.Errorf(
			"failed to get workouts from repository: %w",
			err,
		)
	}

	return workouts, nil
}

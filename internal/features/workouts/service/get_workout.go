package workouts_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

func (s *WorkoutsService) GetWorkout(
	ctx context.Context,
	userID uuid.UUID,
	workoutID uuid.UUID,
) (domain.Workout, error) {
	workout, err := s.workoutsRepository.GetWorkout(ctx, userID, workoutID)
	if err != nil {
		return domain.Workout{}, fmt.Errorf("get workout from repository: %w", err)
	}

	return workout, nil
}

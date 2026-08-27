package workouts_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *WorkoutsService) DeleteWorkout(
	ctx context.Context,
	userID uuid.UUID,
	workoutID uuid.UUID,
) error {
	if err := s.workoutsRepository.DeleteWorkout(ctx, userID, workoutID); err != nil {
		return fmt.Errorf("delete workout from repository: %w", err)
	}

	return nil
}

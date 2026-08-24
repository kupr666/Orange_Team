package workouts_service

import (
	"context"
	"fmt"
)

func (s *WorkoutsService) DeleteWorkout(
	ctx context.Context,
	workoutID int,
) error {
    if err := s.workoutsRepository.DeleteWorkout(ctx, workoutID); err != nil {
		return fmt.Errorf("delete workout from repository: %w", err)
	}

	return nil
}
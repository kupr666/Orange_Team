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
		return fmt.Errorf(
			"delete workout from repository: %w",
			err,
		)
	}

	if err := s.userRepository.UpdateUserWorkoutScore(ctx, userID); err != nil {
		return fmt.Errorf(
			"update user workout score after delete: %w",
			err,
		)
	}

	return nil
}

package habits_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
)

func (s *HabitsService) DeleteHabit(ctx context.Context, userID uuid.UUID, habitID uuid.UUID) error {
	if userID == uuid.Nil || habitID == uuid.Nil {
		return fmt.Errorf("user ID and habit ID are required: %w", core_errors.ErrInvalidArgument)
	}

	if err := s.habitsRepository.DeleteHabit(ctx, userID, habitID); err != nil {
		return fmt.Errorf("delete habit from repository: %w", err)
	}
	return nil
}

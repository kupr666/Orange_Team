package habits_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
)

func (s *HabitsService) CompleteHabit(
	ctx context.Context,
	userID uuid.UUID,
	habitID uuid.UUID,
) (domain.Habit, error) {
	if userID == uuid.Nil || habitID == uuid.Nil {
		return domain.Habit{}, fmt.Errorf("user ID and habit ID are required: %w", core_errors.ErrInvalidArgument)
	}

	habit, err := s.habitsRepository.GetHabit(ctx, userID, habitID)
	if err != nil {
		return domain.Habit{}, fmt.Errorf("get habit from repository: %w", err)
	}

	today := s.now().UTC()
	completed, changed := habit.Complete(today)
	if !changed {
		return completed, nil
	}

	updated, err := s.habitsRepository.CompleteHabit(ctx, userID, completed)
	if err != nil {
		return domain.Habit{}, fmt.Errorf("save habit completion: %w", err)
	}

	return updated.ViewAt(today), nil
}

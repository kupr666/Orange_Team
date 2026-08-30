package habits_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
)

func (s *HabitsService) GetHabits(ctx context.Context, userID uuid.UUID) ([]domain.Habit, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("user ID is empty: %w", core_errors.ErrInvalidArgument)
	}

	habits, err := s.habitsRepository.GetHabits(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get habits from repository: %w", err)
	}

	today := s.today()
	for i := range habits {
		habits[i] = habits[i].ViewAt(today)
	}
	return habits, nil
}

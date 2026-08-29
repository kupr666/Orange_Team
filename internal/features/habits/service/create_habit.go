package habits_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

func (s *HabitsService) CreateHabit(
	ctx context.Context,
	userID uuid.UUID,
	name string,
	description string,
) (domain.Habit, error) {
	habit := domain.Habit{UserID: userID, Name: name, Description: description}
	if err := habit.ValidateForCreation(); err != nil {
		return domain.Habit{}, fmt.Errorf("validate habit: %w", err)
	}

	created, err := s.habitsRepository.CreateHabit(ctx, habit)
	if err != nil {
		return domain.Habit{}, fmt.Errorf("create habit in repository: %w", err)
	}

	return created.ViewAt(s.now()), nil
}

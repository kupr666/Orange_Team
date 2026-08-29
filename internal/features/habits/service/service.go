package habits_service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

type HabitsRepository interface {
	CreateHabit(ctx context.Context, habit domain.Habit) (domain.Habit, error)
	GetHabits(ctx context.Context, userID uuid.UUID) ([]domain.Habit, error)
	GetHabit(ctx context.Context, userID uuid.UUID, habitID uuid.UUID) (domain.Habit, error)
	CompleteHabit(ctx context.Context, userID uuid.UUID, habit domain.Habit) (domain.Habit, error)
	DeleteHabit(ctx context.Context, userID uuid.UUID, habitID uuid.UUID) error
}

type HabitsService struct {
	habitsRepository HabitsRepository
	now              func() time.Time
}

func NewHabitsService(habitsRepository HabitsRepository) *HabitsService {
	return &HabitsService{
		habitsRepository: habitsRepository,
		now:              time.Now,
	}
}

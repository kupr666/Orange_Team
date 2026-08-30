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
	location         *time.Location
	now              func() time.Time
}

func NewHabitsService(
	habitsRepository HabitsRepository,
	location *time.Location,
) *HabitsService {
	return &HabitsService{
		habitsRepository: habitsRepository,
		location:         location,
		now:              time.Now,
	}
}

func (s *HabitsService) today() time.Time {
	localNow := s.now().In(s.location)

	return time.Date(
		localNow.Year(),
		localNow.Month(),
		localNow.Day(),
		0,
		0,
		0,
		0,
		time.UTC,
	)
}

package habits_postgres_repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"
)

type HabitModel struct {
	ID                uuid.UUID
	Version           int64
	UserID            uuid.UUID
	Name              string
	Description       string
	CurrentStreak     int
	LastCompletedDate *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (m *HabitModel) Scan(row core_postgres_pool.Row) error {
	return row.Scan(
		&m.ID,
		&m.Version,
		&m.UserID,
		&m.Name,
		&m.Description,
		&m.CurrentStreak,
		&m.LastCompletedDate,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
}

func domainFromModel(model HabitModel) domain.Habit {
	return domain.NewHabit(
		model.ID,
		model.Version,
		model.UserID,
		model.Name,
		model.Description,
		model.CurrentStreak,
		model.LastCompletedDate,
		model.CreatedAt,
		model.UpdatedAt,
	)
}

func domainsFromModels(models []HabitModel) []domain.Habit {
	habits := make([]domain.Habit, len(models))
	for i, model := range models {
		habits[i] = domainFromModel(model)
	}
	return habits
}

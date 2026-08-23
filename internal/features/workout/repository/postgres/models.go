package workouts_postgres_repository

import (
	"time"

	"github.com/google/uuid"

	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"
)

type WorkoutModel struct {
	ID                       uuid.UUID
	Version                  int
	UserID                   uuid.UUID
	Status                   string
	StartedAt                *time.Time
	CompletedAt              *time.Time
	CreatedAt                time.Time
	UpdatedAt                time.Time
	WorkoutScore             int
	Intensity                *int
	PersonalScoreCoefficient int
}

func (m *WorkoutModel) Scan(row core_postgres_pool.Row) error {
	return row.Scan(
		&m.ID,
		&m.Version,
		&m.UserID,
		&m.Status,
		&m.StartedAt,
		&m.CompletedAt,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.WorkoutScore,
		&m.Intensity,
		&m.PersonalScoreCoefficient,
	)
}

func domainFromModel(model WorkoutModel) domain.Workout {
	return domain.NewWorkout(
		model.ID,
		model.Version,
		model.UserID,
		model.Status,
		model.StartedAt,
		model.CompletedAt,
		model.CreatedAt,
		model.UpdatedAt,
		model.WorkoutScore,
		model.Intensity,
		model.PersonalScoreCoefficient,
	)
}

func domainsFromModels(models []WorkoutModel) []domain.Workout {
	domains := make([]domain.Workout, len(models))

	for i, model := range models {
		domains[i] = domainFromModel(model)
	}

	return domains
}

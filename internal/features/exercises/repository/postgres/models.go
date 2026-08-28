package exercises_postgres_repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"
)

type ExerciseModel struct {
	ID          uuid.UUID
	Version     int
	Name        string
	Description string
	Difficulty  int
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	Type        string
}

func domainFromModel(model ExerciseModel) domain.Exercise {
	return domain.NewExercise(
		model.ID,
		model.Version,
		model.Name,
		model.Description,
		model.Difficulty,
		model.CreatedAt,
		model.UpdatedAt,
		model.Type,
	)
}

func domainsFromModels(models []ExerciseModel) []domain.Exercise {
	domains := make([]domain.Exercise, len(models))

	for i, model := range models {
		domains[i] = domainFromModel(model)
	}

	return domains
}

func (m *ExerciseModel) Scan(row core_postgres_pool.Row) error {
	return row.Scan(
		&m.ID,
		&m.Version,
		&m.Name,
		&m.Description,
		&m.Difficulty,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.Type,
	)
}

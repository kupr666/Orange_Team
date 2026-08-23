package exercises_postgres_repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

type ExerciseModel struct {
	ID          uuid.UUID
	Version     int64
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

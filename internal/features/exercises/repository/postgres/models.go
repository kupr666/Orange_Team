package exercises_postgres_repository

import (
	"time"

	"github.com/kupr666/Orange_Team/internal/core/domain"
)

type ExerciseModel struct {
	ID          int
	Version     int64
	Name        string
	Description string
	Difficulty  int
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	Type        string
}

func exerciseDomainFromModel(exerciseModel ExerciseModel) domain.Exercise {
	return domain.NewExercise(
		exerciseModel.ID,
		exerciseModel.Version,
		exerciseModel.Name,
		exerciseModel.Description,
		exerciseModel.Difficulty,
		exerciseModel.CreatedAt,
		exerciseModel.UpdatedAt,
		exerciseModel.Type,
	)
}

func exerciseDomainsFromModels(exercises []ExerciseModel) []domain.Exercise {
	exercisesDomains := make([]domain.Exercise, len(exercises))

	for i, model := range exercises {
		exercisesDomains[i] = exerciseDomainFromModel(model)
	}

	return exercisesDomains
}

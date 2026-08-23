package exercises_service

import (
	"context"

	"github.com/kupr666/Orange_Team/internal/core/domain"
)

type ExercisesService struct {
	exercisesRepository ExercisesRepository
}

type ExercisesRepository interface {
	GetExercises(
		ctx context.Context,
	) ([]domain.Exercise, error)
}

func NewExercisesService(repo ExercisesRepository) *ExercisesService {
	return &ExercisesService{
		exercisesRepository: repo,
	}
}

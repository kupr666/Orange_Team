package exercises_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

type ExercisesService struct {
	exercisesRepository ExercisesRepository
}

type ExercisesRepository interface {
	GetExercise(
		ctx context.Context,
		exerciseID uuid.UUID,
	) (domain.Exercise, error)

	GetExercises(
		ctx context.Context,
	) ([]domain.Exercise, error)

	CreateExercise(
		ctx context.Context,
		exercise domain.Exercise,
	) (domain.Exercise, error)
}

func NewExercisesService(repo ExercisesRepository) *ExercisesService {
	return &ExercisesService{
		exercisesRepository: repo,
	}
}

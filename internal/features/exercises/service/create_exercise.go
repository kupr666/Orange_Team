package exercises_service

import (
	"context"

	"github.com/kupr666/Orange_Team/internal/core/domain"
)

func (s *ExercisesService) CreateExercise(
	ctx context.Context,
	exercise domain.Exercise,
) (domain.Exercise, error) {
	// 1. exercise.Validate()
	// 2. newExercise := repo.Save(exercise)
	// 3. return newExercise
}



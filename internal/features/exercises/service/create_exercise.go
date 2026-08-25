package exercises_service

import (
	"context"

	"github.com/kupr666/Orange_Team/internal/core/domain"
)

func (s *ExercisesService) CreateExercise(
	ctx context.Context,
	name string,
	description string,
	difficulty int,
	exerciseType string,
) (domain.Exercise, error) {
	if err := exercise.Validate(); err != nil {
		return domain.Exercise{}, err
	}
	newExercise, errL := s.CreateExercise()

	// 1. exercise.Validate()
	// 2. newExercise := repo.Save(exercise)
	// 3. return newExercise
}

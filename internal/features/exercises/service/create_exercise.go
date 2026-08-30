package exercises_service

import (
	"context"
	"fmt"

	"github.com/kupr666/Orange_Team/internal/core/domain"
)

func (s *ExercisesService) CreateExercise(
	ctx context.Context,
	name string,
	description string,
	difficulty int,
	exerciseType string,
) (domain.Exercise, error) {
	exercise := domain.Exercise{
		Name:        name,
		Description: description,
		Difficulty:  difficulty,
		Type:        exerciseType,
	}

	if err := exercise.Validate(); err != nil {
		return domain.Exercise{}, fmt.Errorf(
			"validate exercise: %w",
			err,
		)
	}

	createdExercise, err := s.exercisesRepository.CreateExercise(ctx, exercise)
	if err != nil {
		return domain.Exercise{}, fmt.Errorf(
			"create exercise in repository: %w",
			err,
		)
	}

	return createdExercise, nil
}

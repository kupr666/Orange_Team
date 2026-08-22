package exercises_service

import (
	"context"
	"fmt"

	"github.com/kupr666/Orange_Team/internal/core/domain"
)
func (s *ExercisesService) GetExercises(
	ctx context.Context,
) ([]domain.Exercise, error) {
	
	exercises, err := s.exercisesRepository.GetExercises(ctx)
	if err != nil {
		return []domain.Exercise{}, fmt.Errorf(
			"failed to get tasks from repository: %w",
			err,
		)
	}

	return exercises, nil
}



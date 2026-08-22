package exercises_postgres_repository

import (
	"context"

	"github.com/kupr666/Orange_Team/internal/core/domain"
)

func (r *ExercisesRepository) GetExercises(
	ctx context.Context,
) ([]domain.Exercise, error) {

	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	
}
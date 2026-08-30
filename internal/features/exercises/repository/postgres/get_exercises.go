package exercises_postgres_repository

import (
	"context"
	"fmt"

	"github.com/kupr666/Orange_Team/internal/core/domain"
)

func (r *ExercisesRepository) GetExercises(
	ctx context.Context,
) ([]domain.Exercise, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT
			id,
			version,
			name,
			description,
			difficulty,
			created_at,
			updated_at,
			type
		FROM app.exercises
		ORDER BY id ASC;
	`

	rows, err := r.pool.Query(ctx, query);
	if err != nil {
		return nil, fmt.Errorf(
			"select exercises: %w",
			err,
		)
	}
	defer rows.Close()

	var exercisesModels []ExerciseModel
	for rows.Next() {
		var exerciseModel ExerciseModel
		if err := exerciseModel.Scan(rows); err != nil {
			return nil, fmt.Errorf(
				"scan exercises: %w",
				err,
			)
		}

		exercisesModels = append(exercisesModels, exerciseModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"next rows: %w",
			err,
		)
	}

	exercisesDomains := domainsFromModels(exercisesModels)

	return exercisesDomains, nil
}

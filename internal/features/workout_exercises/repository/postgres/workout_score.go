package workout_exercises_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"
)

func (r *WorkoutExercisesRepository) UpdateWorkoutScore(
	ctx context.Context,
	workoutID uuid.UUID,
	score int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE app.workouts
		SET
			workout_score = $2,
			version = version + 1,
			updated_at = NOW()
		WHERE id = $1
	`

	cmdTag, err := r.pool.Exec(ctx, query, workoutID, score)
	if err != nil {
		return fmt.Errorf("update workout score: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf(
			"workout with id='%s': %w",
			workoutID,
			core_errors.ErrNotFound,
		)
	}

	return nil
}

func (r *WorkoutExercisesRepository) GetPersonalScoreCoefficient(
	ctx context.Context,
	workoutID uuid.UUID,
) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT personal_score_coefficient
		FROM app.workouts
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, workoutID)

	var coefficient int
	if err := row.Scan(&coefficient); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return 0, fmt.Errorf(
				"workout with id='%s': %w",
				workoutID,
				core_errors.ErrNotFound,
			)
		}

		return 0, fmt.Errorf("scan personal score coefficient: %w", err)
	}

	return coefficient, nil
}

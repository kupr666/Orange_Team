package workouts_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
)

func (r *WorkoutsRepository) DeleteWorkout(
	ctx context.Context,
	workoutID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM app.workouts
	WHERE id=$1;
	`

	cmdTag, err := r.pool.Exec(ctx, query, workoutID)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
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

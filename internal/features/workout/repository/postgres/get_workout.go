package workouts_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"
)

func (r *WorkoutsRepository) GetWorkout(
	ctx context.Context,
	workoutID uuid.UUID,
) (domain.Workout, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
    SELECT
        id,
        version,
        user_id,
        status,
        started_at,
        completed_at,
        created_at,
        updated_at,
        workout_score,
        intensity,
        personal_score_coefficient
    FROM app.workouts
    WHERE id = $1;
	`

	row := r.pool.QueryRow(ctx, query, workoutID)

	var workoutModel WorkoutModel
	if err := workoutModel.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Workout{}, fmt.Errorf(
				"workout with id='%s': %w",
				workoutID,
				core_errors.ErrNotFound,
			)
		}

		return domain.Workout{}, fmt.Errorf("scan error: %w", err)
	}

	workoutDomain := domainFromModel(workoutModel)

	return workoutDomain, nil
}

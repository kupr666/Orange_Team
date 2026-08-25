package workouts_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"
)

func (r *WorkoutsRepository) CreateWorkout(
	ctx context.Context,
	userID uuid.UUID,
) (domain.Workout, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	workoutID := uuid.New()
	now := time.Now().UTC()

	// TODO: calculate the personal score coefficient when its business logic is defined.
	const personalScoreCoefficient = 1

	query := `
		INSERT INTO app.workouts (
			id,
			user_id,
			status,
			created_at,
			updated_at,
			personal_score_coefficient
		)
		SELECT
			$1,
			users.id,
			'planned',
			$3,
			$3,
			$4
		FROM app.users AS users
		WHERE users.id = $2
		RETURNING
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
			personal_score_coefficient;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		workoutID,
		userID,
		now,
		personalScoreCoefficient,
	)

	var model WorkoutModel
	if err := model.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Workout{}, fmt.Errorf(
				"user with id='%s': %w",
				userID,
				core_errors.ErrNotFound,
			)
		}

		return domain.Workout{}, fmt.Errorf(
			"scan created workout: %w",
			err,
		)
	}

	return domainFromModel(model), nil
}

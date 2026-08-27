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

func (r *WorkoutsRepository) PatchWorkout(
	ctx context.Context,
	userID uuid.UUID,
	workout domain.Workout,
) (domain.Workout, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE app.workouts
	SET
		version = version + 1,
		status = $2,
		started_at = $3,
		completed_at = $4,
		updated_at = NOW(),
		workout_score = $5,
		intensity = $6,
		personal_score_coefficient = $7
	WHERE id = $1 AND user_id = $8 AND version = $9 
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
		personal_score_coefficient
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		workout.ID,
		workout.Status,
		workout.StartedAt,
		workout.CompletedAt,
		workout.WorkoutScore,
		workout.Intensity,
		workout.PersonalScoreCoefficient,
		workout.UserID,
		workout.Version,
	)

	var workoutModel WorkoutModel
	err := row.Scan(
		&workoutModel.ID,
		&workoutModel.Version,
		&workoutModel.UserID,
		&workoutModel.Status,
		&workoutModel.StartedAt,
		&workoutModel.CompletedAt,
		&workoutModel.CreatedAt,
		&workoutModel.UpdatedAt,
		&workoutModel.WorkoutScore,
		&workoutModel.Intensity,
		&workoutModel.PersonalScoreCoefficient,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Workout{}, fmt.Errorf(
				"workout with id='%s' concurrently updated: %w",
				workout.ID,
				core_errors.ErrConflict,
			)
		}
		return domain.Workout{}, fmt.Errorf("update workout: %w", err)
	}

	workoutDomain := domainFromModel(workoutModel)

	return workoutDomain, nil
}

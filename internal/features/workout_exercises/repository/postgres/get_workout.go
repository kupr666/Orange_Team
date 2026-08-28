package workout_exercises_postgres_repository

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

type workoutModel struct {
	ID                       uuid.UUID
	Version                  int
	UserID                   uuid.UUID
	Status                   string
	StartedAt                *time.Time
	CompletedAt              *time.Time
	CreatedAt                time.Time
	UpdatedAt                time.Time
	WorkoutScore             int
	Intensity                *int
	PersonalScoreCoefficient int
}

func (r *WorkoutExercisesRepository) GetWorkout(
	ctx context.Context,
	userID uuid.UUID,
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
		WHERE user_id = $1 AND id = $2
	`

	row := r.pool.QueryRow(ctx, query, userID, workoutID)

	var model workoutModel
	if err := row.Scan(
		&model.ID,
		&model.Version,
		&model.UserID,
		&model.Status,
		&model.StartedAt,
		&model.CompletedAt,
		&model.CreatedAt,
		&model.UpdatedAt,
		&model.WorkoutScore,
		&model.Intensity,
		&model.PersonalScoreCoefficient,
	); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Workout{}, fmt.Errorf(
				"workout with id='%s': %w",
				workoutID,
				core_errors.ErrNotFound,
			)
		}

		return domain.Workout{}, fmt.Errorf("scan workout: %w", err)
	}

	return domain.NewWorkout(
		model.ID,
		model.Version,
		model.UserID,
		model.Status,
		model.StartedAt,
		model.CompletedAt,
		model.CreatedAt,
		model.UpdatedAt,
		model.WorkoutScore,
		model.Intensity,
		model.PersonalScoreCoefficient,
	), nil
}

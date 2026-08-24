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

func (r *WorkoutsRepository) CreateExercise(
	ctx context.Context,
	userID uuid.UUID,
	exercise domain.CreatedExercise,
) (domain.CreatedExercise, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO app.created_exercises (
			id,
			exercise_id,
			workout_id,
			weight,
			sets,
			reps,
			duration,
			created_at,
			completed,
			exercise_load
		)
		SELECT
			$1,
			$2,
			workouts.id,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11
		FROM app.workouts AS workouts
		WHERE workouts.id = $3
			AND workouts.user_id = $4
			AND EXISTS (
				SELECT 1
				FROM app.exercises
				WHERE exercises.id = $2
			)
		RETURNING
			id,
			version,
			exercise_id,
			workout_id,
			weight,
			sets,
			reps,
			duration,
			created_at,
			updated_at,
			completed,
			exercise_load;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		exercise.ID,
		exercise.ExerciseID,
		exercise.WorkoutID,
		userID,
		exercise.Weight,
		exercise.Sets,
		exercise.Reps,
		exercise.Duration,
		exercise.CreatedAt,
		exercise.Completed,
		exercise.ExerciseLoad,
	)

	var model CreatedExerciseModel
	if err := model.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.CreatedExercise{}, fmt.Errorf(
				"workout or exercise not found: %w",
				core_errors.ErrNotFound,
			)
		}

		return domain.CreatedExercise{}, fmt.Errorf(
			"scan created exercise: %w",
			err,
		)
	}

	return createdExerciseDomainFromModel(model), nil
}

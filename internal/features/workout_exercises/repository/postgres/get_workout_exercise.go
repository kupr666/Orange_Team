package workout_exercises_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"
)

func (r *WorkoutExercisesRepository) GetWorkoutExercise(
	ctx context.Context,
	workoutID uuid.UUID,
	workoutExerciseID uuid.UUID,
) (domain.WorkoutExercise, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT
		id,
		version,
		workout_id,
		exercise_id,
		weight,
		sets,
		reps,
		duration,
		completed,
		exercise_load,
		created_at,
		updated_at
	FROM app.workout_exercises
	WHERE workout_id = $1 AND id = $2
	`

	row := r.pool.QueryRow(ctx, query, workoutID, workoutExerciseID)

	var workoutExerciseModel WorkoutExerciseModel
	if err := workoutExerciseModel.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.WorkoutExercise{}, fmt.Errorf(
				"workout exercise with id='%s' in workout='%s': %w",
				workoutExerciseID,
				workoutID,
				core_errors.ErrNotFound,
			)
		}
		return domain.WorkoutExercise{}, fmt.Errorf(
			"scan workout exercise: %w",
			err,
		)
	}

	workoutExerciseDomain := domainFromModel(workoutExerciseModel)

	return workoutExerciseDomain, nil
}

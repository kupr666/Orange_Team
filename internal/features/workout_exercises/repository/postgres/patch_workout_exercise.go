package workout_exercises_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"
)

func (r *WorkoutExercisesRepository) PatchWorkoutExercise(
	ctx context.Context,
	workoutExercise domain.WorkoutExercise,
) (domain.WorkoutExercise, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE app.workout_exercises
		SET
			weight = $4,
			sets = $5,
			reps = $6,
			duration = $7,
			completed = $8,
			exercise_load = $9,
			version = version + 1,
			updated_at = NOW()
		WHERE id = $1 AND workout_id = $2 AND version = $3
		RETURNING id, version, workout_id, exercise_id, weight, sets, reps, duration,
			completed, exercise_load, created_at, updated_at
		`

	row := r.pool.QueryRow(
		ctx,
		query,
		workoutExercise.ID,
		workoutExercise.WorkoutID,
		workoutExercise.Version,
		workoutExercise.Weight,
		workoutExercise.Sets,
		workoutExercise.Reps,
		workoutExercise.Duration,
		workoutExercise.Completed,
		workoutExercise.ExerciseLoad,
	)

	var workoutExerciseModel WorkoutExerciseModel
	if err := workoutExerciseModel.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.WorkoutExercise{}, fmt.Errorf(
				"workout exercise with id='%s' concurrently accessed: %w",
				workoutExercise.ID,
				core_errors.ErrConflict,
			)
		}
		return domain.WorkoutExercise{}, fmt.Errorf(
			"scan error: %w",
			err,
		)
	}

	workoutExerciseDomain := domainFromModel(workoutExerciseModel)

	return workoutExerciseDomain, nil
}

package workout_exercises_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"
)

func (r *WorkoutExercisesRepository) CreateWorkoutExercise(
	ctx context.Context,
	workoutExercise domain.WorkoutExercise,
) (domain.WorkoutExercise, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO app.workout_exercises (
		id, version, workout_id, exercise_id, weight, sets, reps, duration,
		completed, exercise_load, created_at, updated_at
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	RETURNING id, version, workout_id, exercise_id, weight, sets, reps, duration,
		completed, exercise_load, created_at, updated_at
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		workoutExercise.ID,
		workoutExercise.Version,
		workoutExercise.WorkoutID,
		workoutExercise.ExerciseID,
		workoutExercise.Weight,
		workoutExercise.Sets,
		workoutExercise.Reps,
		workoutExercise.Duration,
		workoutExercise.Completed,
		workoutExercise.ExerciseLoad,
		workoutExercise.CreatedAt,
		workoutExercise.UpdatedAt,
	)

	var workoutExerciseModel WorkoutExerciseModel
	if err := workoutExerciseModel.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return domain.WorkoutExercise{}, fmt.Errorf(
				"%v: exercise '%s' or workout '%s' not found: %w",
				err,
				workoutExercise.ExerciseID,
				workoutExercise.WorkoutID,
				core_errors.ErrNotFound,
			)
		}

		return domain.WorkoutExercise{}, fmt.Errorf(
			"scan created workout exercise: %w",
			err,
		)
	}

	workoutExerciseDomain := domainFromModel(workoutExerciseModel)

	return workoutExerciseDomain, nil
}

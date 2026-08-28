package workout_exercises_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

func (r *WorkoutExercisesRepository) GetWorkoutExercises(
	ctx context.Context,
	workoutID uuid.UUID,
) ([]domain.WorkoutExercise, error) {
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
	WHERE workout_id = $1
	ORDER BY created_at ASC
	`

	rows, err := r.pool.Query(ctx, query, workoutID)
	if err != nil {
		return nil, fmt.Errorf(
			"select workout exercises: %w",
			err,
		)
	}
	defer rows.Close()

	var workoutExerciseModels []WorkoutExerciseModel

	for rows.Next() {
		var workoutExerciseModel WorkoutExerciseModel
		if err := workoutExerciseModel.Scan(rows); err != nil {
			return nil, fmt.Errorf(
				"scan workout exercise: %w",
				err,
			)
		}
		workoutExerciseModels = append(workoutExerciseModels, workoutExerciseModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"next rows: %w",
			err,
		)
	}

	workoutDomains := domainFromModels(workoutExerciseModels)

	return workoutDomains, nil
}

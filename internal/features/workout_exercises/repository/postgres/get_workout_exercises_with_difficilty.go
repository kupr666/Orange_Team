package workout_exercises_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

func (r *WorkoutExercisesRepository) GetWorkoutExercisesWithDifficulty(
	ctx context.Context,
	workoutID uuid.UUID,
) ([]domain.WorkoutExerciseWithDifficulty, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT
		app.workout_exercises.id,
		app.workout_exercises.version,
		app.workout_exercises.workout_id,
		app.workout_exercises.exercise_id,
		app.workout_exercises.weight,
		app.workout_exercises.sets,
		app.workout_exercises.reps,
		app.workout_exercises.duration,
		app.workout_exercises.completed,
		app.workout_exercises.exercise_load,
		app.workout_exercises.created_at,
		app.workout_exercises.updated_at,
		app.exercises.difficulty
	FROM app.workout_exercises
	JOIN app.exercises ON app.workout_exercises.exercise_id = app.exercises.id
	WHERE app.workout_exercises.workout_id = $1
	ORDER BY app.workout_exercises.created_at ASC, app.workout_exercises.id ASC
	`

	rows, err := r.pool.Query(ctx, query, workoutID)
	if err != nil {
		return nil, fmt.Errorf(
			"select workout exercises with difficulty: %w",
			err,
		)
	}
	defer rows.Close()

	var workoutExerciseWithDifficulty []domain.WorkoutExerciseWithDifficulty
	for rows.Next() {
		var workoutExerciseModel WorkoutExerciseModel
		var difficulty int
		if err := rows.Scan(
			&workoutExerciseModel.ID,
			&workoutExerciseModel.Version,
			&workoutExerciseModel.WorkoutID,
			&workoutExerciseModel.ExerciseID,
			&workoutExerciseModel.Weight,
			&workoutExerciseModel.Sets,
			&workoutExerciseModel.Reps,
			&workoutExerciseModel.Duration,
			&workoutExerciseModel.Completed,
			&workoutExerciseModel.ExerciseLoad,
			&workoutExerciseModel.CreatedAt,
			&workoutExerciseModel.UpdatedAt,
			&difficulty,
		); err != nil {
			return nil, fmt.Errorf(
				"scan workout exercise with difficulty: %w",
				err,
			)
		}
		workoutExerciseWithDifficulty = append(workoutExerciseWithDifficulty, domain.WorkoutExerciseWithDifficulty{
			WorkoutExercise: domainFromModel(workoutExerciseModel),
			Difficulty:      difficulty,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"rows iteration: %w",
			err,
		)
	}
	return workoutExerciseWithDifficulty, nil
}

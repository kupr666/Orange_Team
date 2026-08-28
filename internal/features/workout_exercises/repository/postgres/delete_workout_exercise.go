package workout_exercises_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
)

func (r *WorkoutExercisesRepository) DeleteWorkoutExercise(
	ctx context.Context,
	workoutID uuid.UUID,
	workoutExerciseID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		DELETE FROM app.workout_exercises
		WHERE id = $1 AND workout_id = $2
		`

	cmdTag, err := r.pool.Exec(ctx, query, workoutExerciseID, workoutID)
	if err != nil {
		return fmt.Errorf(
			"delete workout exercise: %w",
			err,
		)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf(
			"workout exercise with id='%s' in workout='%s': %w",
			workoutExerciseID,
			workoutID,
			core_errors.ErrNotFound,
		)
	}
	return nil
}

package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *UsersRepository) UpdateUserWorkoutScore(
	ctx context.Context,
	userID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE app.users
		SET user_workout_score = (
			SELECT COALESCE(SUM(workout_score), 0)
			FROM app.workouts
			WHERE user_id = $1 AND status = 'completed'
		),
		version = version + 1,
		updated_at = NOW()
		WHERE id = $1
	`

	cmdTag, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf(
			"update user workout score: %w",
			err,
		)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf(
			"user with id='%s' not found",
			userID,
		)
	}
	return nil
}

package habits_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
)

func (r *HabitsRepository) DeleteHabit(ctx context.Context, userID uuid.UUID, habitID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	commandTag, err := r.pool.Exec(
		ctx,
		`DELETE FROM app.habits WHERE user_id = $1 AND id = $2;`,
		userID,
		habitID,
	)
	if err != nil {
		return fmt.Errorf("delete habit: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("habit with id='%s': %w", habitID, core_errors.ErrNotFound)
	}

	return nil
}

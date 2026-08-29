package habits_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"
)

func (r *HabitsRepository) CompleteHabit(
	ctx context.Context,
	userID uuid.UUID,
	habit domain.Habit,
) (domain.Habit, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE app.habits
		SET version = version + 1,
			current_streak = $3,
			last_completed_date = $4,
			updated_at = NOW()
		WHERE user_id = $1 AND id = $2 AND version = $5
		RETURNING id, version, user_id, name, description, current_streak,
			last_completed_date, created_at, updated_at;
	`
	row := r.pool.QueryRow(
		ctx,
		query,
		userID,
		habit.ID,
		habit.CurrentStreak,
		habit.LastCompletedDate,
		habit.Version,
	)

	var model HabitModel
	if err := model.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Habit{}, fmt.Errorf(
				"habit with id='%s' concurrently updated: %w",
				habit.ID,
				core_errors.ErrConflict,
			)
		}
		return domain.Habit{}, fmt.Errorf("update habit completion: %w", err)
	}

	return domainFromModel(model), nil
}

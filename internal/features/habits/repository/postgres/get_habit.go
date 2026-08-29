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

func (r *HabitsRepository) GetHabit(
	ctx context.Context,
	userID uuid.UUID,
	habitID uuid.UUID,
) (domain.Habit, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, user_id, name, description, current_streak,
			last_completed_date, created_at, updated_at
		FROM app.habits
		WHERE user_id = $1 AND id = $2;
	`
	row := r.pool.QueryRow(ctx, query, userID, habitID)

	var model HabitModel
	if err := model.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Habit{}, fmt.Errorf("habit with id='%s': %w", habitID, core_errors.ErrNotFound)
		}
		return domain.Habit{}, fmt.Errorf("scan habit: %w", err)
	}

	return domainFromModel(model), nil
}

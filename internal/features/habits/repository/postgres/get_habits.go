package habits_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

func (r *HabitsRepository) GetHabits(ctx context.Context, userID uuid.UUID) ([]domain.Habit, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, user_id, name, description, current_streak,
			last_completed_date, created_at, updated_at
		FROM app.habits
		WHERE user_id = $1
		ORDER BY created_at DESC, id ASC;
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("select habits: %w", err)
	}
	defer rows.Close()

	var models []HabitModel
	for rows.Next() {
		var model HabitModel
		if err := model.Scan(rows); err != nil {
			return nil, fmt.Errorf("scan habit: %w", err)
		}
		models = append(models, model)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate habits: %w", err)
	}

	return domainsFromModels(models), nil
}

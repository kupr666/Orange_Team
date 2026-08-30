package habits_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"
)

func (r *HabitsRepository) CreateHabit(ctx context.Context, habit domain.Habit) (domain.Habit, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	now := time.Now().UTC()
	query := `
		INSERT INTO app.habits (id, user_id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $5)
		RETURNING id, version, user_id, name, description, current_streak,
			last_completed_date, created_at, updated_at;
	`

	row := r.pool.QueryRow(ctx, query, uuid.New(), habit.UserID, habit.Name, habit.Description, now)
	var model HabitModel
	if err := model.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesUnique) {
			return domain.Habit{}, fmt.Errorf(
				"habit with name=%q already exists for user: %w",
				habit.Name,
				core_errors.ErrConflict,
			)
		}
		return domain.Habit{}, fmt.Errorf("scan created habit: %w", err)
	}

	return domainFromModel(model), nil
}

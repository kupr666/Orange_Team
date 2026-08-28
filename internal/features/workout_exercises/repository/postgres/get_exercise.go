package workout_exercises_postgres_repository

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

type exerciseModel struct {
	ID          uuid.UUID
	Version     int
	Name        string
	Description string
	Difficulty  int
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	Type        string
}

func (r *WorkoutExercisesRepository) GetExercise(
	ctx context.Context,
	exerciseID uuid.UUID,
) (domain.Exercise, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT
			id,
			version,
			name,
			description,
			difficulty,
			created_at,
			updated_at,
			type
		FROM app.exercises
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, exerciseID)

	var model exerciseModel
	if err := row.Scan(
		&model.ID,
		&model.Version,
		&model.Name,
		&model.Description,
		&model.Difficulty,
		&model.CreatedAt,
		&model.UpdatedAt,
		&model.Type,
	); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Exercise{}, fmt.Errorf(
				"exercise with id='%s': %w",
				exerciseID,
				core_errors.ErrNotFound,
			)
		}

		return domain.Exercise{}, fmt.Errorf("scan exercise: %w", err)
	}

	return domain.NewExercise(
		model.ID,
		model.Version,
		model.Name,
		model.Description,
		model.Difficulty,
		model.CreatedAt,
		model.UpdatedAt,
		model.Type,
	), nil
}

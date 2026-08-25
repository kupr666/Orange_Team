package exercises_postgres_repository

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

func (r *ExercisesRepository) CreateExercise(
	ctx context.Context,
	exercise domain.Exercise,
) (domain.Exercise, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	exerciseID := uuid.New()
	createdAt := time.Now().UTC()

	query := `
		INSERT INTO app.exercises (
			id,
			name,
			description,
			difficulty,
			created_at,
			type
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
			id,
			version,
			name,
			description,
			difficulty,
			created_at,
			updated_at,
			type;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		exerciseID,
		exercise.Name,
		exercise.Description,
		exercise.Difficulty,
		createdAt,
		exercise.Type,
	)

	var model ExerciseModel
	err := row.Scan(
		&model.ID,
		&model.Version,
		&model.Name,
		&model.Description,
		&model.Difficulty,
		&model.CreatedAt,
		&model.UpdatedAt,
		&model.Type,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesUnique) {
			return domain.Exercise{}, fmt.Errorf(
				"exercise with name=%q already exists: %w",
				exercise.Name,
				core_errors.ErrConflict,
			)
		}

		return domain.Exercise{}, fmt.Errorf(
			"scan created exercise: %w",
			err,
		)
	}

	return domainFromModel(model), nil
}

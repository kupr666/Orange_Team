package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) GetUser(
	ctx context.Context,
	userID uuid.UUID,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT
			id,
			version,
			mail,
			full_name,
			created_at,
			updated_at,
			user_workout_score,
			sex,
			weight_grams,
			birth_date,
			height_cm
		FROM app.users
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, userID)

	var userModel UserModel
	if err := userModel.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with id='%s': %w",
				userID,
				core_errors.ErrNotFound,
			)
		}
		return domain.User{}, fmt.Errorf("scan user: %w", err)
	}

	return domainFromModel(userModel), nil
}

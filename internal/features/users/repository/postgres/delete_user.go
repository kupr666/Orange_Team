package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
)

func (r *UsersRepository) DeleteUser(
	ctx context.Context,
	userID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELITE FROM app.users
    WHERE id=$1;
    `

	cmdTag, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf(
			"exec query: %w",
			err,
		)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf(
			"user with id='%d': %w",
			userID,
			core_errors.ErrNotFound,
		)
	}

	return nil
}

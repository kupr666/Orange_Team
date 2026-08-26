package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) PatchUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE app.users
	SET
		version = version + 1,
		sex = $2,
		weight_grams = $3,
		birth_date = $4,
		height_cm = $5,
		updated_at = NOW()
	WHERE id = $1 AND version = $6
	RETURNING
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
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.ID,
		user.Sex,
		user.WeightGrams,
		user.BirthDate,
		user.HeightCM,
		user.Version,
	)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.Email,
		&userModel.FullName,
		&userModel.CreatedAt,
		&userModel.UpdatedAt,
		&userModel.UserWorkoutScore,
		&userModel.Sex,
		&userModel.WeightGrams,
		&userModel.BirthDate,
		&userModel.HeightCM,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with id='%s' concurrently updated: %w",
				user.ID,
				core_errors.ErrConflict,
			)
		}
		return domain.User{}, fmt.Errorf("update user: %w", err)
	}

	userDomain := domainFromModel(userModel)

	return userDomain, nil
}

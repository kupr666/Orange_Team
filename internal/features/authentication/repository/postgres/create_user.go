package authentication_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"
)

func (r *AuthenticationRepository) CreateUser(
	ctx context.Context,
	email string,
	passwordHash string,
	fullName string,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO app.users (mail, pass_hash, full_name)
	VALUES ($1, $2, $3)
	RETURNING id, version, mail, full_name, created_at, updated_at,
	user_workout_score, sex, weight_grams, birth_date, height_cm;
	`

	row := r.pool.QueryRow(ctx, query, email, passwordHash, fullName)

	var userModel UserModel

	err := row.Scan(
		&userModel.id,
		&userModel.version,
		&userModel.mail,
		&userModel.fullName,
		&userModel.createdAt,
		&userModel.updatedAt,
		&userModel.userWorkoutScore,
		&userModel.sex,
		&userModel.weightGrams,
		&userModel.birthDate,
		&userModel.heightCm,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesUnique) {
			return domain.User{}, fmt.Errorf(
				"user with email %q already exists: %w",
				email,
				core_errors.ErrConflict,
			)
		}

		return domain.User{}, fmt.Errorf("scan created user: %w", err)
	}

	domainUser := domainUserFromModel(userModel)

	return domainUser, nil
}

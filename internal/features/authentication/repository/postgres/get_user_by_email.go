package authentication_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"
	authentication_domain "github.com/kupr666/Orange_Team/internal/features/authentication/domain"
)

func (r *AuthenticationRepository) GetLoginCredentialsByEmail(
	ctx context.Context,
	email string,
) (authentication_domain.StoredCredentials, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT
		id,
		pass_hash,
		role
	FROM app.users
	WHERE LOWER(mail) = LOWER($1)
	`
	var credentials CredentialsModel
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&credentials.id,
		&credentials.passwordHash,
		&credentials.role,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return authentication_domain.StoredCredentials{}, fmt.Errorf(
				"credentials for email: %q: %w",
				email,
				core_errors.ErrNotFound,
			)
		}

		return authentication_domain.StoredCredentials{}, fmt.Errorf(
			"scan login credentials: %w",
			err,
		)
	}

	credentialsDomain := storedCredentialsFromModel(credentials)

	return credentialsDomain, nil
}

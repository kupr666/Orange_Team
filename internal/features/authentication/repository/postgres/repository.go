package authentication_postgres_repository

import core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"

type AuthenticationRepository struct {
	pool core_postgres_pool.Pool
}

func NewAuthenticationRepository(pool core_postgres_pool.Pool) *AuthenticationRepository {
	return &AuthenticationRepository{
		pool: pool,
	}
}

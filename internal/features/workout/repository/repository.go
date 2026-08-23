package workout_repository

import core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"

type WorkoutsRepository struct {
	pool core_postgres_pool.Pool
}

func NewWorkoutsRepository(pool core_postgres_pool.Pool) *WorkoutsRepository {
	return &WorkoutsRepository{
		pool: pool,
	}
}

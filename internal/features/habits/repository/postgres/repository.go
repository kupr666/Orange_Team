package habits_postgres_repository

import core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"

type HabitsRepository struct {
	pool core_postgres_pool.Pool
}

func NewHabitsRepository(pool core_postgres_pool.Pool) *HabitsRepository {
	return &HabitsRepository{pool: pool}
}

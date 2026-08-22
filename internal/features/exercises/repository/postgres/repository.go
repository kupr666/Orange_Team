package exercises_postgres_repository

import core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"

type ExercisesRepository struct {
	pool core_postgres_pool.Pool
}

func NewExercisesRepository(pool core_postgres_pool.Pool) *ExercisesRepository {
	return &ExercisesRepository{
		pool: pool,
	}
}

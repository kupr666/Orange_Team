package workout_exercises_postgres_repository

import core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"

type WorkoutExercisesRepository struct {
	pool core_postgres_pool.Pool
}

func NewWorkoutExercisesRepository(pool core_postgres_pool.Pool) *WorkoutExercisesRepository {
	return &WorkoutExercisesRepository{
		pool: pool,
	}
}

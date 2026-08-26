package leaderboard_postgres_repository

import core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"

type LeaderboardRepository struct {
	pool core_postgres_pool.Pool
}

func NewLeaderboardRepository(
	pool core_postgres_pool.Pool,
) *LeaderboardRepository {
	return &LeaderboardRepository{pool: pool}
}

package leaderboard_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

func (s *LeaderboardService) GetWeeklyLeaderboard(
	ctx context.Context,
	userID uuid.UUID,
	limit int,
) (domain.Leaderboard, error) {
	return s.getPublishedLeaderboard(
		ctx,
		userID,
		limit,
		domain.LeaderboardPeriodWeekly,
	)
}

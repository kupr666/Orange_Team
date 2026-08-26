package leaderboard_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

func (s *LeaderboardService) GetDailyLeaderboard(
	ctx context.Context,
	userID uuid.UUID,
	limit int,
) (domain.Leaderboard, error) {
	if err := validateArguments(userID, limit); err != nil {
		return domain.Leaderboard{}, err
	}

	now := s.now()
	window := dailyWindow(now, s.location)
	ranking, err := s.leaderboardRepository.GetLiveLeaderboard(
		ctx,
		userID,
		window.start,
		window.end,
		limit,
	)
	if err != nil {
		return domain.Leaderboard{}, fmt.Errorf(
			"get live daily leaderboard: %w",
			err,
		)
	}

	return domain.Leaderboard{
		Period:        domain.LeaderboardPeriodDaily,
		Status:        domain.LeaderboardStatusLive,
		PeriodStart:   window.start,
		PeriodEnd:     window.end,
		Timezone:      s.location.String(),
		Metric:        domain.LeaderboardScoringRule,
		GeneratedAt:   now,
		NextRefreshAt: window.nextRefresh,
		Ranking:       ranking,
	}, nil
}

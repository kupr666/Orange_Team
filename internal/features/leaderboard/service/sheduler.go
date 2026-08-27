package leaderboard_service

import (
	"context"
	"time"

	"github.com/kupr666/Orange_Team/internal/core/domain"
)

func (s *LeaderboardService) FinalizeClosedPeriods(ctx context.Context) error {
	now := s.now()
	periods := []string{
		domain.LeaderboardPeriodWeekly,
		domain.LeaderboardPeriodMonthly,
	}
	for _, p := range periods {
		window, err := previousWindow(now, p, s.location)
		if err != nil {
			continue
		}
		_ = s.repo.CreateSnapshot(ctx, p, window.start, window.end, s.location.String())
	}
	return nil
}

func (s *LeaderboardService) RunScheduler(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	_ = s.FinalizeClosedPeriods(ctx)

	for {
		select {
		case <-ticker.C:
			_ = s.FinalizeClosedPeriods(ctx)
		case <-ctx.Done():
			return
		}
	}
}

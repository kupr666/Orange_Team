package leaderboard_service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

func (s *LeaderboardService) GetDaily(
	ctx context.Context,
	userID uuid.UUID,
	limit int,
) (domain.Leaderboard, error) {
	now := s.now()
	window := dailyWindow(now, s.location)

	ranking, err := s.repo.GetLive(ctx, userID, window.start, window.end, limit)
	if err != nil {
		return domain.Leaderboard{}, err
	}

	return domain.Leaderboard{
		Period:        domain.LeaderboardPeriodDaily,
		Status:        domain.LeaderboardStatusLive,
		PeriodStart:   window.start,
		PeriodEnd:     window.end,
		Timezone:      s.location.String(),
		Metric:        domain.LeaderboardScoringRule,
		GeneratedAt:   now,
		NextRefreshAt: window.end,
		Ranking:       ranking,
	}, nil
}

func (s *LeaderboardService) GetSnapshot(
	ctx context.Context,
	userID uuid.UUID,
	periodType string,
	limit int,
) (domain.Leaderboard, error) {
	now := s.now()
	window, err := previousWindow(now, periodType, s.location)
	if err != nil {
		return domain.Leaderboard{}, err
	}

	snapshot, err := s.repo.GetSnapshot(
		ctx,
		userID,
		periodType,
		window.start,
		limit)
	if err != nil {
		return domain.Leaderboard{}, err
	}

	return domain.Leaderboard{
		Period:        periodType,
		Status:        domain.LeaderboardStatusPublished,
		PeriodStart:   snapshot.PeriodStart,
		PeriodEnd:     snapshot.PeriodEnd,
		Timezone:      snapshot.Timezone,
		Metric:        domain.LeaderboardScoringRule,
		GeneratedAt:   snapshot.PublishedAt,
		NextRefreshAt: window.nextRefresh,
		Ranking:       snapshot.Ranking,
	}, nil
}

type periodWindow struct {
	start       time.Time
	end         time.Time
	nextRefresh time.Time
}

func dailyWindow(now time.Time, loc *time.Location) periodWindow {
	localNow := now.In(loc)
	start := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 0, 0, 0, loc)
	end := start.AddDate(0, 0, 1)
	return periodWindow{start: start, end: end, nextRefresh: end}
}

func previousWindow(now time.Time, periodType string, loc *time.Location) (periodWindow, error) {
	switch periodType {
	case domain.LeaderboardPeriodWeekly:
		return previousWeeklyWindow(now, loc), nil
	case domain.LeaderboardPeriodMonthly:
		return previousMonthlyWindow(now, loc), nil
	default:
		return periodWindow{}, fmt.Errorf(
			"unsupported period %q",
			periodType,
		)
	}
}

func previousWeeklyWindow(now time.Time, loc *time.Location) periodWindow {
	localNow := now.In(loc)
	today := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 0, 0, 0, loc)
	daysSinceMonday := (int(today.Weekday()) + 6) % 7
	currentWeekStart := today.AddDate(0, 0, -daysSinceMonday)
	previousWeekStart := currentWeekStart.AddDate(0, 0, -7)
	return periodWindow{
		start:       previousWeekStart,
		end:         currentWeekStart,
		nextRefresh: currentWeekStart.AddDate(0, 0, 7),
	}
}

func previousMonthlyWindow(now time.Time, loc *time.Location) periodWindow {
	localNow := now.In(loc)
	currentMonthStart := time.Date(localNow.Year(), localNow.Month(), 1, 0, 0, 0, 0, loc)
	previousMonthStart := currentMonthStart.AddDate(0, -1, 0)
	return periodWindow{
		start:       previousMonthStart,
		end:         currentMonthStart,
		nextRefresh: currentMonthStart.AddDate(0, 1, 0),
	}
}

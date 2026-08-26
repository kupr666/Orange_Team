package leaderboard_service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
)

func (s *LeaderboardService) getPublishedLeaderboard(
	ctx context.Context,
	userID uuid.UUID,
	limit int,
	periodType string,
) (domain.Leaderboard, error) {
	if err := validateArguments(userID, limit); err != nil {
		return domain.Leaderboard{}, err
	}

	now := s.now()
	window, err := s.closedWindow(now, periodType)
	if err != nil {
		return domain.Leaderboard{}, err
	}

	if err := s.leaderboardRepository.CreateSnapshot(
		ctx,
		periodType,
		window.start,
		window.end,
		s.location.String(),
	); err != nil {
		return domain.Leaderboard{}, fmt.Errorf(
			"ensure %s leaderboard snapshot: %w",
			periodType,
			err,
		)
	}

	snapshot, err := s.leaderboardRepository.GetSnapshotLeaderboard(
		ctx,
		userID,
		periodType,
		window.start,
		limit,
	)
	if err != nil {
		return domain.Leaderboard{}, fmt.Errorf(
			"get %s leaderboard snapshot: %w",
			periodType,
			err,
		)
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

func validateArguments(userID uuid.UUID, limit int) error {
	if userID == uuid.Nil {
		return fmt.Errorf("user ID is empty: %w", core_errors.ErrInvalidArgument)
	}
	if limit < 1 || limit > MaximumLimit {
		return fmt.Errorf(
			"limit must be between 1 and %d: %w",
			MaximumLimit,
			core_errors.ErrInvalidArgument,
		)
	}
	return nil
}

func (s *LeaderboardService) closedWindow(
	now time.Time,
	periodType string,
) (periodWindow, error) {
	switch periodType {
	case domain.LeaderboardPeriodWeekly:
		return previousWeeklyWindow(now, s.location), nil
	case domain.LeaderboardPeriodMonthly:
		return previousMonthlyWindow(now, s.location), nil
	default:
		return periodWindow{}, fmt.Errorf(
			"unsupported leaderboard period %q: %w",
			periodType,
			core_errors.ErrInvalidArgument,
		)
	}
}

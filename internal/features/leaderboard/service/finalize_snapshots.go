package leaderboard_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kupr666/Orange_Team/internal/core/domain"
)

func (s *LeaderboardService) FinalizeClosedPeriods(
	ctx context.Context,
	now time.Time,
) error {
	periodTypes := []string{
		domain.LeaderboardPeriodWeekly,
		domain.LeaderboardPeriodMonthly,
	}

	var snapshotErrors []error
	for _, periodType := range periodTypes {
		window, err := s.closedWindow(now, periodType)
		if err != nil {
			snapshotErrors = append(snapshotErrors, err)
			continue
		}

		if err := s.leaderboardRepository.CreateSnapshot(
			ctx,
			periodType,
			window.start,
			window.end,
			s.location.String(),
		); err != nil {
			snapshotErrors = append(
				snapshotErrors,
				fmt.Errorf("finalize %s snapshot: %w", periodType, err),
			)
		}
	}

	return errors.Join(snapshotErrors...)
}

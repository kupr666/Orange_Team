package leaderboard_transport_http

import (
	"fmt"
	"net/http"

	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_http_request "github.com/kupr666/Orange_Team/internal/core/transport/http/request"
	leaderboard_service "github.com/kupr666/Orange_Team/internal/features/leaderboard/service"
)

func leaderboardRequestParams(r *http.Request) (string, int, error) {
	period := r.URL.Query().Get("period")
	if period == "" {
		period = domain.LeaderboardPeriodDaily
	}

	allowed := map[string]bool{
		domain.LeaderboardPeriodDaily:   true,
		domain.LeaderboardPeriodWeekly:  true,
		domain.LeaderboardPeriodMonthly: true,
	}
	if !allowed[period] {
		return "", 0, fmt.Errorf(
			"period must be one of: daily, weekly, monthly: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	limit := leaderboard_service.DefaultLimit
	if limitParam, err := core_http_request.GetIntQueryParam(r, "limit"); err == nil && limitParam != nil {
		limit = *limitParam
	}
	if limit < 1 || limit > leaderboard_service.MaximumLimit {
		return "", 0, fmt.Errorf(
			"limit must be between 1 and %d: %w",
			leaderboard_service.MaximumLimit,
			core_errors.ErrInvalidArgument,
		)
	}

	return period, limit, nil
}

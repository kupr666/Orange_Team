package leaderboard_transport_http

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_http_request "github.com/kupr666/Orange_Team/internal/core/transport/http/request"
	leaderboard_service "github.com/kupr666/Orange_Team/internal/features/leaderboard/service"
)

func leaderboardRequestParams(r *http.Request) (uuid.UUID, int, error) {
	userID, err := core_http_request.GetUUIDQueryParam(r, "user_id")
	if err != nil {
		return uuid.Nil, 0, fmt.Errorf("get user_id query parameter: %w", err)
	}
	if userID == nil || *userID == uuid.Nil {
		return uuid.Nil, 0, fmt.Errorf(
			"user_id query parameter is required: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	limit := leaderboard_service.DefaultLimit
	limitParam, err := core_http_request.GetIntQueryParam(r, "limit")
	if err != nil {
		return uuid.Nil, 0, fmt.Errorf("get limit query parameter: %w", err)
	}
	if limitParam != nil {
		limit = *limitParam
	}
	if limit < 1 || limit > leaderboard_service.MaximumLimit {
		return uuid.Nil, 0, fmt.Errorf(
			"limit must be between 1 and %d: %w",
			leaderboard_service.MaximumLimit,
			core_errors.ErrInvalidArgument,
		)
	}

	return *userID, limit, nil
}

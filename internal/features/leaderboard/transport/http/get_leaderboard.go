package leaderboard_transport_http

import (
	"net/http"

	core_auth "github.com/kupr666/Orange_Team/internal/core/auth"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

func (h *LeaderboardHTTPHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	principal, ok := core_auth.PrincipalFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(
			core_errors.ErrUnauthorized,
			"authenticated user is missing,",
		)
		return
	}
	userID := principal.UserID

	period, limit, err := leaderboardRequestParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"invalid leaderboard query parameters",
		)
		return
	}

	var board domain.Leaderboard
	switch period {
	case domain.LeaderboardPeriodDaily:
		board, err = h.leaderboardService.GetDaily(ctx, userID, limit)
	case domain.LeaderboardPeriodWeekly, domain.LeaderboardPeriodMonthly:
		board, err = h.leaderboardService.GetSnapshot(ctx, userID, period, limit)
	}
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get leaderboard",
		)
		return
	}

	w.Header().Set("Cache-Control", "private, no-store")
	responseHandler.JSONResponse(leaderboardResponseFromDomain(board), http.StatusOK)
}

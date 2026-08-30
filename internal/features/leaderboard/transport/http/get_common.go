package leaderboard_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

type leaderboardGetter func(
	ctx context.Context,
	userID uuid.UUID,
	limit int,
) (domain.Leaderboard, error)

func (h *LeaderboardHTTPHandler) getLeaderboard(
	w http.ResponseWriter,
	r *http.Request,
	get leaderboardGetter,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID, limit, err := leaderboardRequestParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid leaderboard query parameters")
		return
	}

	leaderboard, err := get(ctx, userID, limit)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get leaderboard")
		return
	}

	w.Header().Set("Cache-Control", "private, no-store")
	responseHandler.JSONResponse(
		leaderboardResponseFromDomain(leaderboard),
		http.StatusOK,
	)
}

package leaderboard_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_http_server "github.com/kupr666/Orange_Team/internal/core/transport/http/server"
)

type LeaderboardService interface {
	GetDaily(
		ctx context.Context,
		userID uuid.UUID,
		limit int,
	) (domain.Leaderboard, error)
	GetSnapshot(
		ctx context.Context,
		userID uuid.UUID,
		period string,
		limit int,
	) (domain.Leaderboard, error)
}

type LeaderboardHTTPHandler struct {
	leaderboardService LeaderboardService
}

func NewLeaderboardHTTPHandler(service LeaderboardService) *LeaderboardHTTPHandler {
	return &LeaderboardHTTPHandler{leaderboardService: service}
}

func (h *LeaderboardHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/leaderboard",
			Handler: h.GetLeaderboard,
		},
	}
}

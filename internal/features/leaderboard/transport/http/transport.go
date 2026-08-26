package leaderboard_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_http_server "github.com/kupr666/Orange_Team/internal/core/transport/http/server"
)

type LeaderboardService interface {
	GetDailyLeaderboard(
		ctx context.Context,
		userID uuid.UUID,
		limit int,
	) (domain.Leaderboard, error)

	GetWeeklyLeaderboard(
		ctx context.Context,
		userID uuid.UUID,
		limit int,
	) (domain.Leaderboard, error)

	GetMonthlyLeaderboard(
		ctx context.Context,
		userID uuid.UUID,
		limit int,
	) (domain.Leaderboard, error)
}

type LeaderboardHTTPHandler struct {
	leaderboardService LeaderboardService
}

func NewLeaderboardHTTPHandler(
	leaderboardService LeaderboardService,
) *LeaderboardHTTPHandler {
	return &LeaderboardHTTPHandler{leaderboardService: leaderboardService}
}

func (h *LeaderboardHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/leaderboard/daily",
			Handler: h.GetDailyLeaderboard,
		},
		{
			Method:  http.MethodGet,
			Path:    "/leaderboard/weekly",
			Handler: h.GetWeeklyLeaderboard,
		},
		{
			Method:  http.MethodGet,
			Path:    "/leaderboard/monthly",
			Handler: h.GetMonthlyLeaderboard,
		},
	}
}

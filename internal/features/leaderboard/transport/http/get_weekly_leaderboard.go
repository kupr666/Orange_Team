package leaderboard_transport_http

import "net/http"

func (h *LeaderboardHTTPHandler) GetWeeklyLeaderboard(
	w http.ResponseWriter,
	r *http.Request,
) {
	h.getLeaderboard(w, r, h.leaderboardService.GetWeeklyLeaderboard)
}

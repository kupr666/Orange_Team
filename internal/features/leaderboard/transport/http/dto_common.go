package leaderboard_transport_http

import (
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

type LeaderboardUserResponse struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
}

type LeaderboardEntryResponse struct {
	Rank           *int64                  `json:"rank"`
	User           LeaderboardUserResponse `json:"user"`
	Score          int64                   `json:"score"`
	LastActivityAt *time.Time              `json:"last_activity_at"`
	IsCurrentUser  bool                    `json:"is_current_user"`
}

type CurrentUserLeaderboardResponse struct {
	Rank           *int64                  `json:"rank"`
	User           LeaderboardUserResponse `json:"user"`
	Score          int64                   `json:"score"`
	LastActivityAt *time.Time              `json:"last_activity_at"`
	Eligible       bool                    `json:"eligible"`
	InTop          bool                    `json:"in_top"`
}

type LeaderboardResponse struct {
	Period        string                         `json:"period"`
	Status        string                         `json:"status"`
	PeriodStart   time.Time                      `json:"period_start"`
	PeriodEnd     time.Time                      `json:"period_end"`
	Timezone      string                         `json:"timezone"`
	Metric        string                         `json:"metric"`
	GeneratedAt   time.Time                      `json:"generated_at"`
	NextRefreshAt time.Time                      `json:"next_refresh_at"`
	Items         []LeaderboardEntryResponse     `json:"items"`
	CurrentUser   CurrentUserLeaderboardResponse `json:"current_user"`
}

func leaderboardResponseFromDomain(leaderboard domain.Leaderboard) LeaderboardResponse {
	items := make([]LeaderboardEntryResponse, len(leaderboard.Ranking.Entries))
	for i, entry := range leaderboard.Ranking.Entries {
		items[i] = LeaderboardEntryResponse{
			Rank: entry.Rank,
			User: LeaderboardUserResponse{
				ID:       entry.UserID,
				FullName: entry.FullName,
			},
			Score:          entry.Score,
			LastActivityAt: entry.LastActivityAt,
			IsCurrentUser:  entry.IsCurrentUser,
		}
	}

	currentUser := leaderboard.Ranking.CurrentUser
	return LeaderboardResponse{
		Period:        leaderboard.Period,
		Status:        leaderboard.Status,
		PeriodStart:   leaderboard.PeriodStart,
		PeriodEnd:     leaderboard.PeriodEnd,
		Timezone:      leaderboard.Timezone,
		Metric:        leaderboard.Metric,
		GeneratedAt:   leaderboard.GeneratedAt,
		NextRefreshAt: leaderboard.NextRefreshAt,
		Items:         items,
		CurrentUser: CurrentUserLeaderboardResponse{
			Rank: currentUser.Rank,
			User: LeaderboardUserResponse{
				ID:       currentUser.UserID,
				FullName: currentUser.FullName,
			},
			Score:          currentUser.Score,
			LastActivityAt: currentUser.LastActivityAt,
			Eligible:       currentUser.Rank != nil,
			InTop:          currentUser.IsInTop,
		},
	}
}

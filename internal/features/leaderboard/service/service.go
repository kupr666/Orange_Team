package leaderboard_service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

const (
	DefaultLimit = 50
	MaximumLimit = 100
)

type LeaderboardRepository interface {
	GetLive(ctx context.Context, userID uuid.UUID, periodStart, periodEnd time.Time, limit int) (domain.LeaderboardRanking, error)
	GetSnapshot(ctx context.Context, userID uuid.UUID, periodType string, periodStart time.Time, limit int) (domain.LeaderboardSnapshot, error)
	CreateSnapshot(ctx context.Context, periodType string, periodStart, periodEnd time.Time, timezone string) error
}

type LeaderboardService struct {
	repo     LeaderboardRepository
	location *time.Location
	now      func() time.Time
}

func NewLeaderboardService(repo LeaderboardRepository, location *time.Location) *LeaderboardService {
	return &LeaderboardService{
		repo:     repo,
		location: location,
		now:      time.Now,
	}
}

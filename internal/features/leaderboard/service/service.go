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
	GetLiveLeaderboard(
		ctx context.Context,
		userID uuid.UUID,
		periodStart time.Time,
		periodEnd time.Time,
		limit int,
	) (domain.LeaderboardRanking, error)

	CreateSnapshot(
		ctx context.Context,
		periodType string,
		periodStart time.Time,
		periodEnd time.Time,
		timezone string,
	) error

	GetSnapshotLeaderboard(
		ctx context.Context,
		userID uuid.UUID,
		periodType string,
		periodStart time.Time,
		limit int,
	) (domain.LeaderboardSnapshot, error)
}

type LeaderboardService struct {
	leaderboardRepository LeaderboardRepository
	location              *time.Location
	now                   func() time.Time
}

func NewLeaderboardService(
	repository LeaderboardRepository,
	location *time.Location,
) *LeaderboardService {
	return &LeaderboardService{
		leaderboardRepository: repository,
		location:              location,
		now:                   time.Now,
	}
}

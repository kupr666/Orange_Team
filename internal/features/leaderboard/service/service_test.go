package leaderboard_service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

type leaderboardRepositoryStub struct {
	liveStart      time.Time
	liveEnd        time.Time
	liveLimit      int
	createdPeriods []string
	createdStarts  []time.Time
	snapshotPeriod string
	snapshotStart  time.Time
	snapshotLimit  int
	ranking        domain.LeaderboardRanking
}

func (r *leaderboardRepositoryStub) GetLiveLeaderboard(
	_ context.Context,
	_ uuid.UUID,
	periodStart time.Time,
	periodEnd time.Time,
	limit int,
) (domain.LeaderboardRanking, error) {
	r.liveStart = periodStart
	r.liveEnd = periodEnd
	r.liveLimit = limit
	return r.ranking, nil
}

func (r *leaderboardRepositoryStub) CreateSnapshot(
	_ context.Context,
	periodType string,
	periodStart time.Time,
	_ time.Time,
	_ string,
) error {
	r.createdPeriods = append(r.createdPeriods, periodType)
	r.createdStarts = append(r.createdStarts, periodStart)
	return nil
}

func (r *leaderboardRepositoryStub) GetSnapshotLeaderboard(
	_ context.Context,
	_ uuid.UUID,
	periodType string,
	periodStart time.Time,
	limit int,
) (domain.LeaderboardSnapshot, error) {
	r.snapshotPeriod = periodType
	r.snapshotStart = periodStart
	r.snapshotLimit = limit
	return domain.LeaderboardSnapshot{
		PeriodStart: periodStart,
		PeriodEnd:   periodStart.AddDate(0, 0, 7),
		Timezone:    "Europe/Moscow",
		PublishedAt: periodStart.AddDate(0, 0, 7).Add(time.Minute),
		Ranking:     r.ranking,
	}, nil
}

func TestGetDailyLeaderboardUsesCurrentDay(t *testing.T) {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		t.Fatalf("load location: %v", err)
	}

	userID := uuid.New()
	repository := &leaderboardRepositoryStub{
		ranking: domain.LeaderboardRanking{
			CurrentUser: domain.LeaderboardEntry{UserID: userID},
		},
	}
	service := NewLeaderboardService(repository, location)
	service.now = func() time.Time {
		return time.Date(2026, time.August, 26, 14, 0, 0, 0, location)
	}

	result, err := service.GetDailyLeaderboard(context.Background(), userID, 25)
	if err != nil {
		t.Fatalf("get daily leaderboard: %v", err)
	}

	wantStart := time.Date(2026, time.August, 26, 0, 0, 0, 0, location)
	wantEnd := time.Date(2026, time.August, 27, 0, 0, 0, 0, location)
	if !repository.liveStart.Equal(wantStart) || !repository.liveEnd.Equal(wantEnd) {
		t.Fatalf(
			"repository window = [%s, %s), want [%s, %s)",
			repository.liveStart,
			repository.liveEnd,
			wantStart,
			wantEnd,
		)
	}
	if repository.liveLimit != 25 {
		t.Fatalf("repository limit = %d, want 25", repository.liveLimit)
	}
	if result.Status != domain.LeaderboardStatusLive {
		t.Fatalf("status = %q, want live", result.Status)
	}
}

func TestGetWeeklyLeaderboardCreatesPreviousWeekSnapshot(t *testing.T) {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		t.Fatalf("load location: %v", err)
	}

	userID := uuid.New()
	repository := &leaderboardRepositoryStub{
		ranking: domain.LeaderboardRanking{
			CurrentUser: domain.LeaderboardEntry{UserID: userID},
		},
	}
	service := NewLeaderboardService(repository, location)
	service.now = func() time.Time {
		return time.Date(2026, time.August, 26, 14, 0, 0, 0, location)
	}

	result, err := service.GetWeeklyLeaderboard(
		context.Background(),
		userID,
		DefaultLimit,
	)
	if err != nil {
		t.Fatalf("get weekly leaderboard: %v", err)
	}

	wantStart := time.Date(2026, time.August, 17, 0, 0, 0, 0, location)
	if len(repository.createdPeriods) != 1 ||
		repository.createdPeriods[0] != domain.LeaderboardPeriodWeekly {
		t.Fatalf("created periods = %v, want [weekly]", repository.createdPeriods)
	}
	if !repository.createdStarts[0].Equal(wantStart) {
		t.Fatalf("snapshot start = %s, want %s", repository.createdStarts[0], wantStart)
	}
	if repository.snapshotPeriod != domain.LeaderboardPeriodWeekly ||
		!repository.snapshotStart.Equal(wantStart) {
		t.Fatalf(
			"read snapshot = (%s, %s), want (weekly, %s)",
			repository.snapshotPeriod,
			repository.snapshotStart,
			wantStart,
		)
	}
	if result.Status != domain.LeaderboardStatusPublished {
		t.Fatalf("status = %q, want published", result.Status)
	}
}

func TestGetMonthlyLeaderboardCreatesPreviousMonthSnapshot(t *testing.T) {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		t.Fatalf("load location: %v", err)
	}

	userID := uuid.New()
	repository := &leaderboardRepositoryStub{
		ranking: domain.LeaderboardRanking{
			CurrentUser: domain.LeaderboardEntry{UserID: userID},
		},
	}
	service := NewLeaderboardService(repository, location)
	service.now = func() time.Time {
		return time.Date(2026, time.August, 26, 14, 0, 0, 0, location)
	}

	result, err := service.GetMonthlyLeaderboard(
		context.Background(),
		userID,
		DefaultLimit,
	)
	if err != nil {
		t.Fatalf("get monthly leaderboard: %v", err)
	}

	wantStart := time.Date(2026, time.July, 1, 0, 0, 0, 0, location)
	if len(repository.createdPeriods) != 1 ||
		repository.createdPeriods[0] != domain.LeaderboardPeriodMonthly {
		t.Fatalf("created periods = %v, want [monthly]", repository.createdPeriods)
	}
	if !repository.createdStarts[0].Equal(wantStart) {
		t.Fatalf("snapshot start = %s, want %s", repository.createdStarts[0], wantStart)
	}
	if result.Period != domain.LeaderboardPeriodMonthly {
		t.Fatalf("period = %q, want monthly", result.Period)
	}
}

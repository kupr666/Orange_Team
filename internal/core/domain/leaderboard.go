package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	LeaderboardPeriodDaily   = "daily"
	LeaderboardPeriodWeekly  = "weekly"
	LeaderboardPeriodMonthly = "monthly"
)

const (
	LeaderboardStatusLive      = "live"
	LeaderboardStatusPublished = "published"
)

const (
	LeaderboardCategoryWorkouts = "workouts"
	LeaderboardScoringRule      = "workouts_count"
)

type Leaderboard struct {
	Period        string
	Status        string
	PeriodStart   time.Time
	PeriodEnd     time.Time
	Timezone      string
	Metric        string
	GeneratedAt   time.Time
	NextRefreshAt time.Time
	Ranking       LeaderboardRanking
}

type LeaderboardRanking struct {
	Entries     []LeaderboardEntry
	CurrentUser LeaderboardEntry
}

type LeaderboardEntry struct {
	Rank           *int64
	UserID         uuid.UUID
	FullName       string
	Score          int64
	LastActivityAt *time.Time
	IsCurrentUser  bool
	IsInTop        bool
}

type LeaderboardSnapshot struct {
	PeriodStart time.Time
	PeriodEnd   time.Time
	Timezone    string
	PublishedAt time.Time
	Ranking     LeaderboardRanking
}

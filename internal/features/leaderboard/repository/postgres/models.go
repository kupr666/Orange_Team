package leaderboard_postgres_repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

type scanner interface {
	Scan(dest ...any) error
}

type LeaderboardEntryModel struct {
	Rank              *int64
	UserID            uuid.UUID
	FullName          string
	Score             int64
	CompletedWorkouts int64
	LastActivityAt    *time.Time
	IsCurrentUser     bool
	IsInTop           bool
}

func (m *LeaderboardEntryModel) Scan(row scanner) error {
	return row.Scan(
		&m.Rank,
		&m.UserID,
		&m.FullName,
		&m.Score,
		&m.CompletedWorkouts,
		&m.LastActivityAt,
		&m.IsCurrentUser,
		&m.IsInTop,
	)
}

func domainFromEntryModel(model LeaderboardEntryModel) domain.LeaderboardEntry {
	return domain.LeaderboardEntry{
		Rank:              model.Rank,
		UserID:            model.UserID,
		FullName:          model.FullName,
		Score:             model.Score,
		CompletedWorkouts: model.CompletedWorkouts,
		LastActivityAt:    model.LastActivityAt,
		IsCurrentUser:     model.IsCurrentUser,
		IsInTop:           model.IsInTop,
	}
}

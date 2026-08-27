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
	Rank           *int
	UserID         uuid.UUID
	FullName       string
	Score          int
	LastActivityAt *time.Time
	IsCurrentUser  bool
	IsInTop        bool
}

func (m *LeaderboardEntryModel) Scan(row scanner) error {
	return row.Scan(
		&m.Rank,
		&m.UserID,
		&m.FullName,
		&m.Score,
		&m.LastActivityAt,
		&m.IsCurrentUser,
		&m.IsInTop,
	)
}

func domainFromEntryModel(model LeaderboardEntryModel) domain.LeaderboardEntry {
	return domain.LeaderboardEntry{
		Rank:           model.Rank,
		UserID:         model.UserID,
		FullName:       model.FullName,
		Score:          model.Score,
		LastActivityAt: model.LastActivityAt,
		IsCurrentUser:  model.IsCurrentUser,
		IsInTop:        model.IsInTop,
	}
}

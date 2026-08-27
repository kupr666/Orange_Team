package leaderboard_postgres_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
)

// -------- GetLive (для daily) --------

func (r *LeaderboardRepository) GetLive(
	ctx context.Context,
	userID uuid.UUID,
	periodStart time.Time,
	periodEnd time.Time,
	limit int,
) (domain.LeaderboardRanking, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	WITH scores AS (
		SELECT
			workouts.user_id,
			COUNT(*)::BIGINT AS score,
			MAX(workouts.completed_at) AS last_activity_at
		FROM app.workouts AS workouts
		WHERE workouts.status = 'completed'
		  AND workouts.completed_at >= $1
		  AND workouts.completed_at < LEAST($2, CURRENT_TIMESTAMP)
		GROUP BY workouts.user_id
	),
	ranked AS (
		SELECT
			scores.*,
			RANK() OVER (ORDER BY scores.score DESC) AS rank
		FROM scores
	),
	top_rows AS (
		SELECT
			ranked.rank,
			ranked.user_id,
			ranked.score,
			ranked.last_activity_at,
			TRUE AS is_in_top,
			ranked.user_id = $4 AS is_current_user
		FROM ranked
		ORDER BY ranked.rank, ranked.last_activity_at, ranked.user_id
		LIMIT $3
	),
	selected AS (
		SELECT * FROM top_rows

		UNION ALL

		SELECT
			ranked.rank,
			ranked.user_id,
			ranked.score,
			ranked.last_activity_at,
			FALSE,
			TRUE
		FROM ranked
		WHERE ranked.user_id = $4
		  AND NOT EXISTS (
			SELECT 1 FROM top_rows WHERE top_rows.user_id = $4
		  )

		UNION ALL

		SELECT
			NULL::BIGINT,
			users.id,
			0::BIGINT,
			NULL::TIMESTAMPTZ,
			FALSE,
			TRUE
		FROM app.users AS users
		WHERE users.id = $4
		  AND NOT EXISTS (
			SELECT 1 FROM ranked WHERE ranked.user_id = $4
		  )
	)
	SELECT
		selected.rank,
		users.id,
		users.full_name,
		selected.score,
		selected.last_activity_at,
		selected.is_current_user,
		selected.is_in_top
	FROM selected
	JOIN app.users AS users ON users.id = selected.user_id
	ORDER BY
		selected.rank NULLS LAST,
		selected.last_activity_at NULLS LAST,
		selected.user_id;
	`

	rows, err := r.pool.Query(ctx, query, periodStart, periodEnd, limit, userID)
	if err != nil {
		return domain.LeaderboardRanking{}, fmt.Errorf("select live leaderboard: %w", err)
	}
	defer rows.Close()

	ranking := domain.LeaderboardRanking{Entries: make([]domain.LeaderboardEntry, 0, limit)}
	currentUserFound := false

	for rows.Next() {
		var model LeaderboardEntryModel
		if err := model.Scan(rows); err != nil {
			return domain.LeaderboardRanking{}, fmt.Errorf("scan live leaderboard entry: %w", err)
		}
		entry := domainFromEntryModel(model)
		if entry.IsInTop {
			ranking.Entries = append(ranking.Entries, entry)
		}
		if entry.IsCurrentUser {
			ranking.CurrentUser = entry
			currentUserFound = true
		}
	}
	if err := rows.Err(); err != nil {
		return domain.LeaderboardRanking{}, fmt.Errorf("iterate live leaderboard entries: %w", err)
	}
	if !currentUserFound {
		return domain.LeaderboardRanking{}, fmt.Errorf("leaderboard user with id='%s': %w", userID, core_errors.ErrNotFound)
	}
	return ranking, nil
}

func (r *LeaderboardRepository) GetSnapshot(
	ctx context.Context,
	userID uuid.UUID,
	periodType string,
	periodStart time.Time,
	limit int,
) (domain.LeaderboardSnapshot, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	WITH snapshot AS (
		SELECT
			snapshots.id,
			snapshots.period_start,
			snapshots.period_end,
			snapshots.timezone,
			snapshots.published_at
		FROM app.leaderboard_snapshots AS snapshots
		WHERE snapshots.category = $1
		  AND snapshots.period_type = $2
		  AND snapshots.period_start = $3
		  AND snapshots.scoring_rule = $4
		LIMIT 1
	),
	top_rows AS (
		SELECT
			entries.rank,
			entries.user_id,
			entries.score,
			entries.last_activity_at,
			TRUE AS is_in_top,
			entries.user_id = $6 AS is_current_user
		FROM app.leaderboard_snapshot_entries AS entries
		JOIN snapshot ON snapshot.id = entries.snapshot_id
		ORDER BY entries.rank, entries.last_activity_at, entries.user_id
		LIMIT $5
	),
	selected AS (
		SELECT * FROM top_rows

		UNION ALL

		SELECT
			entries.rank,
			entries.user_id,
			entries.score,
			entries.last_activity_at,
			FALSE,
			TRUE
		FROM app.leaderboard_snapshot_entries AS entries
		JOIN snapshot ON snapshot.id = entries.snapshot_id
		WHERE entries.user_id = $6
		  AND NOT EXISTS (
			SELECT 1 FROM top_rows WHERE top_rows.user_id = $6
		  )

		UNION ALL

		SELECT
			NULL::BIGINT,
			users.id,
			0::BIGINT,
			NULL::TIMESTAMPTZ,
			FALSE,
			TRUE
		FROM app.users AS users
		CROSS JOIN snapshot
		WHERE users.id = $6
		  AND NOT EXISTS (
			SELECT 1
			FROM app.leaderboard_snapshot_entries AS entries
			WHERE entries.snapshot_id = snapshot.id
			  AND entries.user_id = $6
		  )
	)
	SELECT
		snapshot.period_start,
		snapshot.period_end,
		snapshot.timezone,
		snapshot.published_at,
		selected.rank,
		users.id,
		users.full_name,
		selected.score,
		selected.last_activity_at,
		selected.is_current_user,
		selected.is_in_top
	FROM snapshot
	JOIN selected ON TRUE
	JOIN app.users AS users ON users.id = selected.user_id
	ORDER BY
		selected.rank NULLS LAST,
		selected.last_activity_at NULLS LAST,
		selected.user_id;
	`

	rows, err := r.pool.Query(ctx, query, domain.LeaderboardCategoryWorkouts, periodType, periodStart, domain.LeaderboardScoringRule, limit, userID)
	if err != nil {
		return domain.LeaderboardSnapshot{}, fmt.Errorf(
			"select leaderboard snapshot: %w",
			err,
		)
	}
	defer rows.Close()

	snapshot := domain.LeaderboardSnapshot{
		Ranking: domain.LeaderboardRanking{Entries: make([]domain.LeaderboardEntry, 0, limit)},
	}
	currentUserFound := false
	rowFound := false

	for rows.Next() {
		rowFound = true
		var model LeaderboardEntryModel
		if err := rows.Scan(
			&snapshot.PeriodStart,
			&snapshot.PeriodEnd,
			&snapshot.Timezone,
			&snapshot.PublishedAt,
			&model.Rank,
			&model.UserID,
			&model.FullName,
			&model.Score,
			&model.LastActivityAt,
			&model.IsCurrentUser,
			&model.IsInTop,
		); err != nil {
			return domain.LeaderboardSnapshot{}, fmt.Errorf("scan leaderboard snapshot entry: %w", err)
		}
		entry := domainFromEntryModel(model)
		if entry.IsInTop {
			snapshot.Ranking.Entries = append(snapshot.Ranking.Entries, entry)
		}
		if entry.IsCurrentUser {
			snapshot.Ranking.CurrentUser = entry
			currentUserFound = true
		}
	}
	if err := rows.Err(); err != nil {
		return domain.LeaderboardSnapshot{}, fmt.Errorf(
			"iterate leaderboard snapshot entries: %w",
			err,
		)
	}
	if !rowFound {
		return domain.LeaderboardSnapshot{}, fmt.Errorf(
			"%s leaderboard snapshot starting at %s: %w",
			periodType,
			periodStart.Format(time.RFC3339),
			core_errors.ErrNotFound,
		)
	}
	if !currentUserFound {
		return domain.LeaderboardSnapshot{}, fmt.Errorf(
			"leaderboard user with id='%s': %w",
			userID,
			core_errors.ErrNotFound,
		)
	}
	return snapshot, nil
}

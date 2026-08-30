package leaderboard_postgres_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
)

func (r *LeaderboardRepository) GetLiveLeaderboard(
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
			COALESCE(SUM(workouts.workout_score), 0)::BIGINT AS score,
			COUNT(*)::BIGINT AS completed_workouts,
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
			ranked.completed_workouts,
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
			ranked.completed_workouts,
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
		selected.completed_workouts,
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

	rows, err := r.pool.Query(
		ctx,
		query,
		periodStart,
		periodEnd,
		limit,
		userID,
	)
	if err != nil {
		return domain.LeaderboardRanking{}, fmt.Errorf(
			"select live leaderboard: %w",
			err,
		)
	}
	defer rows.Close()

	ranking := domain.LeaderboardRanking{
		Entries: make([]domain.LeaderboardEntry, 0, limit),
	}
	currentUserFound := false

	for rows.Next() {
		var model LeaderboardEntryModel
		if err := model.Scan(rows); err != nil {
			return domain.LeaderboardRanking{}, fmt.Errorf(
				"scan live leaderboard entry: %w",
				err,
			)
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
		return domain.LeaderboardRanking{}, fmt.Errorf(
			"iterate live leaderboard entries: %w",
			err,
		)
	}

	if !currentUserFound {
		return domain.LeaderboardRanking{}, fmt.Errorf(
			"leaderboard user with id='%s': %w",
			userID,
			core_errors.ErrNotFound,
		)
	}

	return ranking, nil
}

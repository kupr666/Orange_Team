package leaderboard_postgres_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

func (r *LeaderboardRepository) CreateSnapshot(
	ctx context.Context,
	periodType string,
	periodStart time.Time,
	periodEnd time.Time,
	timezone string,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	snapshotID := uuid.New()
	query := `
	WITH new_snapshot AS (
		INSERT INTO app.leaderboard_snapshots (
			id,
			category,
			period_type,
			period_start,
			period_end,
			timezone,
			scoring_rule,
			published_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP)
		ON CONFLICT (
			category,
			period_type,
			period_start,
			scoring_rule
		)
		DO NOTHING
		RETURNING id
	),
	scores AS (
		SELECT
			workouts.user_id,
			COUNT(*)::BIGINT AS score,
			MAX(workouts.completed_at) AS last_activity_at
		FROM app.workouts AS workouts
		WHERE workouts.status = 'completed'
		  AND workouts.completed_at >= $4
		  AND workouts.completed_at < $5
		GROUP BY workouts.user_id
	),
	ranked AS (
		SELECT
			scores.*,
			RANK() OVER (ORDER BY scores.score DESC) AS rank
		FROM scores
	)
	INSERT INTO app.leaderboard_snapshot_entries (
		snapshot_id,
		user_id,
		rank,
		score,
		last_activity_at
	)
	SELECT
		new_snapshot.id,
		ranked.user_id,
		ranked.rank,
		ranked.score,
		ranked.last_activity_at
	FROM new_snapshot
	CROSS JOIN ranked;
	`

	if _, err := r.pool.Exec(
		ctx,
		query,
		snapshotID,
		domain.LeaderboardCategoryWorkouts,
		periodType,
		periodStart,
		periodEnd,
		timezone,
		domain.LeaderboardScoringRule,
	); err != nil {
		return fmt.Errorf(
			"create %s leaderboard snapshot: %w",
			periodType,
			err,
		)
	}
	return nil
}

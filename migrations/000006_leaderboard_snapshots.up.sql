CREATE TABLE IF NOT EXISTS app.leaderboard_snapshots (
    id             UUID         PRIMARY KEY,
    category       VARCHAR(32)  NOT NULL,
    period_type    VARCHAR(16)  NOT NULL,
    period_start   TIMESTAMPTZ  NOT NULL,
    period_end     TIMESTAMPTZ  NOT NULL,
    timezone       VARCHAR(64)  NOT NULL,
    scoring_rule   VARCHAR(64)  NOT NULL,
    published_at   TIMESTAMPTZ  NOT NULL,

    CONSTRAINT leaderboard_snapshots_period_type_check
        CHECK (period_type IN ('weekly', 'monthly')),
    CONSTRAINT leaderboard_snapshots_period_check
        CHECK (period_start < period_end),
    CONSTRAINT leaderboard_snapshots_unique_period
        UNIQUE (category, period_type, period_start, scoring_rule)
);

CREATE TABLE IF NOT EXISTS app.leaderboard_snapshot_entries (
    snapshot_id      UUID         NOT NULL
        REFERENCES app.leaderboard_snapshots (id) ON DELETE CASCADE,
    user_id          UUID         NOT NULL
        REFERENCES app.users (id) ON DELETE CASCADE,
    rank             BIGINT       NOT NULL,
    score            BIGINT       NOT NULL,
    last_activity_at TIMESTAMPTZ,

    PRIMARY KEY (snapshot_id, user_id),
    CONSTRAINT leaderboard_snapshot_entries_rank_check CHECK (rank > 0),
    CONSTRAINT leaderboard_snapshot_entries_score_check CHECK (score >= 0)
);

CREATE INDEX leaderboard_snapshot_entries_rank_idx
    ON app.leaderboard_snapshot_entries (snapshot_id, rank, user_id)
    INCLUDE (score, last_activity_at);

CREATE INDEX workouts_completed_period_user_idx
    ON app.workouts (completed_at, user_id)
    WHERE status = 'completed';
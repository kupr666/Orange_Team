CREATE TABLE app.workout_sessions (
    id           SERIAL      PRIMARY KEY,
    user_id      INTEGER     NOT NULL REFERENCES app.users (id) ON DELETE CASCADE,
    status       VARCHAR(32) NOT NULL DEFAULT 'in_progress',
    started_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at  TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT workout_sessions_status_check
        CHECK (status IN ('in_progress', 'completed', 'cancelled')),
    CONSTRAINT workout_sessions_completion_check
        CHECK (
            (status = 'in_progress' AND finished_at IS NULL)
            OR (status IN ('completed', 'cancelled') AND finished_at IS NOT NULL)
        ),
    CONSTRAINT workout_sessions_time_order_check
        CHECK (finished_at IS NULL OR finished_at >= started_at)
);

CREATE INDEX workout_sessions_user_started_at_idx
    ON app.workout_sessions (user_id, started_at DESC);

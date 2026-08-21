CREATE TABLE IF NOT EXISTS app.workouts (
    id            SERIAL      PRIMARY KEY,
    user_id       INTEGER     NOT NULL REFERENCES app.users (id) ON DELETE CASCADE,
    status        VARCHAR(20) NOT NULL,
    started_at    TIMESTAMPTZ,
    completed_at  TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL,
    updated_at    TIMESTAMPTZ NOT NULL,
    workout_score INTEGER NOT NULL DEFAULT 0
);

ALTER TABLE app.workouts ADD CONSTRAINT workouts_status_check CHECK (status IN ('planned', 'in_progress', 'completed', 'cancelled'));

ALTER TABLE app.workouts ADD CONSTRAINT workouts_started_completed_logic 
    CHECK (
        (status = 'planned' AND started_at IS NULL AND completed_at IS NULL) OR
        (status = 'in_progress' AND started_at IS NOT NULL AND completed_at IS NULL) OR
        (status = 'completed' AND started_at IS NOT NULL AND completed_at IS NOT NULL AND completed_at >= started_at) OR
        (status = 'cancelled' AND completed_at IS NULL)
    );
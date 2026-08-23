CREATE TABLE IF NOT EXISTS app.workouts (
    id                          UUID          PRIMARY KEY,
    version                     BIGINT        NOT NULL DEFAULT 1,
    user_id                     UUID          NOT NULL REFERENCES app.users (id) ON DELETE CASCADE,
    status                      VARCHAR(20)   NOT NULL,
    started_at                  TIMESTAMPTZ,
    completed_at                TIMESTAMPTZ,
    created_at                  TIMESTAMPTZ   NOT NULL,
    updated_at                  TIMESTAMPTZ   NOT NULL,
    workout_score               INTEGER       NOT NULL DEFAULT 0,
    intensity                   SMALLINT,
    personal_score_coefficient  SMALLINT      NOT NULL
);

ALTER TABLE app.workouts ADD CONSTRAINT workouts_status_check CHECK (status IN ('planned', 'in_progress', 'completed', 'cancelled'));
ALTER TABLE app.workouts ADD CONSTRAINT workouts_started_completed_logic 
    CHECK (
        (status = 'planned' AND started_at IS NULL AND completed_at IS NULL) OR
        (status = 'in_progress' AND started_at IS NOT NULL AND completed_at IS NULL) OR
        (status = 'completed' AND started_at IS NOT NULL AND completed_at IS NOT NULL AND completed_at >= started_at) OR
        (status = 'cancelled' AND completed_at IS NULL)
    );

ALTER TABLE app.workouts ADD CONSTRAINT workouts_workout_score_check CHECK (workout_score BETWEEN 0 AND 10000);

ALTER TABLE app.workouts ADD CONSTRAINT workouts_intensity_check CHECK (intensity IS NULL OR intensity BETWEEN 1 AND 10);
ALTER TABLE app.workouts ADD CONSTRAINT workouts_completed_intensity_check CHECK (status <> 'completed' OR intensity IS NOT NULL);

ALTER TABLE app.workouts ADD CONSTRAINT workouts_personal_score_coefficient_check CHECK (personal_score_coefficient BETWEEN 1 AND 10);

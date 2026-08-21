CREATE TABLE app.workout_exercises (
    id                   SERIAL         PRIMARY KEY,
    workout_session_id   INTEGER        NOT NULL REFERENCES app.workout_sessions (id) ON DELETE CASCADE,
    exercise_id          INTEGER        NOT NULL REFERENCES app.exercises (id) ON DELETE RESTRICT,
    sets                 INTEGER        CHECK (sets > 0),
    reps                 INTEGER        CHECK (reps > 0),
    weight               NUMERIC(10, 2) CHECK (weight >= 0),
    position             INTEGER        NOT NULL CHECK (position >= 0),
    duration_seconds     INTEGER        CHECK (duration_seconds > 0),
    points_earned        INTEGER        NOT NULL DEFAULT 0 CHECK (points_earned >= 0),
    created_at           TIMESTAMPTZ    NOT NULL DEFAULT NOW(),

    CONSTRAINT workout_exercises_session_position_unique
        UNIQUE (workout_session_id, position)
);

CREATE INDEX workout_exercises_exercise_id_idx
    ON app.workout_exercises (exercise_id);

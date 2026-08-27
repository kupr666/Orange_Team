CREATE TABLE IF NOT EXISTS app.workout_exercises (
    id          UUID    PRIMARY KEY,
    version     BIGINT  NOT NULL DEFAULT 1,
    exercise_id UUID    NOT NULL REFERENCES app.exercises (id),
    workout_id  UUID    NOT NULL REFERENCES app.workouts (id) ON DELETE CASCADE,
    weight      INTEGER,
    sets        INTEGER,
    reps        INTEGER,
    duration    INTEGER,
    created_at  TIMESTAMPTZ NOT NULL,
    updated_at  TIMESTAMPTZ,
    completed   BOOLEAN NOT NULL DEFAULT false,
    exercise_load INTEGER NOT NULL DEFAULT 0
);

ALTER TABLE app.workout_exercises ADD CONSTRAINT workout_exercises_fields_check 
    CHECK (
        (weight IS NOT NULL AND sets IS NOT NULL AND reps IS NOT NULL AND duration IS NULL) OR
        (weight IS NULL AND sets IS NULL AND reps IS NULL AND duration IS NOT NULL)
    );

ALTER TABLE app.workout_exercises ADD CONSTRAINT workout_exercises_weight_positive  CHECK (weight IS NULL OR weight >= 0);

ALTER TABLE app.workout_exercises ADD CONSTRAINT workout_exercises_sets_positive  CHECK (sets IS NULL OR sets >= 0);

ALTER TABLE app.workout_exercises ADD CONSTRAINT workout_exercises_reps_positive  CHECK (reps IS NULL OR reps >= 0);

ALTER TABLE app.workout_exercises ADD CONSTRAINT workout_exercises_duration_positive  CHECK (duration IS NULL OR duration >= 0);

ALTER TABLE app.workout_exercises ADD CONSTRAINT workout_exercises_exercise_load_check CHECK (exercise_load BETWEEN 0 AND 1000000);

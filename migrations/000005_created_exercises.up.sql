CREATE TABLE IF NOT EXISTS app.created_exercises (
    id          SERIAL PRIMARY KEY,
    exercise_id INTEGER NOT NULL REFERENCES app.exercises (id),
    workout_id  INTEGER NOT NULL REFERENCES app.workouts (id) ON DELETE CASCADE,
    type        VARCHAR(20) NOT NULL,
    weight      INTEGER,
    sets        INTEGER,
    reps        INTEGER,
    duration    INTEGER,
    created_at  TIMESTAMPTZ NOT NULL,
    updated_at  TIMESTAMPTZ,
    completed   BOOLEAN NOT NULL DEFAULT false
);

ALTER TABLE app.created_exercises ADD CONSTRAINT created_exercises_type_check  CHECK (type IN ('weight', 'duration'));

ALTER TABLE app.created_exercises ADD CONSTRAINT created_exercises_fields_check 
    CHECK (
        (type = 'weight' AND 
         weight IS NOT NULL AND sets IS NOT NULL AND reps IS NOT NULL AND duration IS NULL) OR
        (type = 'duration' AND 
         weight IS NULL AND sets IS NULL AND reps IS NULL AND duration IS NOT NULL)
    );

ALTER TABLE app.created_exercises ADD CONSTRAINT created_exercises_weight_positive  CHECK (weight IS NULL OR weight >= 0);

ALTER TABLE app.created_exercises ADD CONSTRAINT created_exercises_sets_positive  CHECK (sets IS NULL OR sets >= 0);

ALTER TABLE app.created_exercises ADD CONSTRAINT created_exercises_reps_positive  CHECK (reps IS NULL OR reps >= 0);

ALTER TABLE app.created_exercises ADD CONSTRAINT created_exercises_duration_positive  CHECK (duration IS NULL OR duration >= 0);
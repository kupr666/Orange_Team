ALTER TABLE app.workout_exercises ADD CONSTRAINT workout_exercises_exercise_load_check CHECK (exercise_load BETWEEN 0 AND 1000000);
DROP TABLE IF EXISTS app.workout_exercises;
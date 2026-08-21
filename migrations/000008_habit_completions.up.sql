CREATE TABLE app.habit_completions (
    id               SERIAL      PRIMARY KEY,
    habit_id         INTEGER     NOT NULL REFERENCES app.habits (id) ON DELETE CASCADE,
    completion_date  DATE        NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL,

    CONSTRAINT habit_completions_habit_date_unique
        UNIQUE (habit_id, completion_date)
);

CREATE TABLE IF NOT EXISTS app.habits (
    id                   UUID         PRIMARY KEY,
    version              BIGINT       NOT NULL DEFAULT 1,
    user_id              UUID         NOT NULL REFERENCES app.users (id) ON DELETE CASCADE,
    name                 VARCHAR(80)  NOT NULL,
    description          VARCHAR(500) NOT NULL DEFAULT '',
    current_streak       INTEGER      NOT NULL DEFAULT 0,
    last_completed_date  DATE,
    created_at           TIMESTAMPTZ  NOT NULL,
    updated_at           TIMESTAMPTZ  NOT NULL
);

ALTER TABLE app.habits ADD CONSTRAINT habits_name_length_check
    CHECK (char_length(trim(name)) BETWEEN 2 AND 80);
ALTER TABLE app.habits ADD CONSTRAINT habits_name_no_leading_trailing_spaces
    CHECK (name = trim(name));
ALTER TABLE app.habits ADD CONSTRAINT habits_description_length_check
    CHECK (char_length(description) <= 500);
ALTER TABLE app.habits ADD CONSTRAINT habits_current_streak_check
    CHECK (current_streak >= 0);

CREATE UNIQUE INDEX idx_habits_unique_user_lower_name
    ON app.habits (user_id, LOWER(name));
CREATE INDEX idx_habits_user_created_at
    ON app.habits (user_id, created_at DESC);

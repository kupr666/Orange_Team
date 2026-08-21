CREATE TABLE app.habits (
    id           SERIAL       PRIMARY KEY,
    user_id      INTEGER      NOT NULL REFERENCES app.users (id) ON DELETE CASCADE,
    name         VARCHAR(100) NOT NULL,
    description  VARCHAR(1000),
    frequency    VARCHAR(32)  NOT NULL,
    is_active    BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMPTZ  NOT NULL,
    updated_at   TIMESTAMPTZ  NOT NULL
);

CREATE INDEX habits_user_id_idx
    ON app.habits (user_id);

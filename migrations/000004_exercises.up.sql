CREATE TABLE app.exercises (
    id           SERIAL       PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    description  VARCHAR(1000),
    base_points  INTEGER      NOT NULL DEFAULT 0 CHECK (base_points >= 0),
    is_active    BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMPTZ  NOT NULL,
    updated_at   TIMESTAMPTZ  NOT NULL
);

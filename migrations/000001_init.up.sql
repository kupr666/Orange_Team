CREATE SCHEMA app;

CREATE TABLE IF NOT EXISTS app.users (
    id          SERIAL PRIMARY KEY,
    mail        VARCHAR(30)  NOT NULL UNIQUE,
    pass_hash   VARCHAR(255) NOT NULL,
    full_name   VARCHAR(50)  NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL,
    updated_at  TIMESTAMPTZ
);
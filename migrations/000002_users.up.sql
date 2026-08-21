CREATE TABLE IF NOT EXISTS app.users (
    id          SERIAL PRIMARY KEY,
    mail        VARCHAR(30)  NOT NULL,
    pass_hash   VARCHAR(255) NOT NULL,
    full_name   VARCHAR(50)  NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL,
    updated_at  TIMESTAMPTZ,
    user_workout_score  INTEGER      NOT NULL DEFAULT 0
);

ALTER TABLE app.users ADD CONSTRAINT users_mail_length_check CHECK (char_length(mail) BETWEEN 5 AND 30);
ALTER TABLE app.users ADD CONSTRAINT users_mail_format_check CHECK (mail ~* '^[a-z0-9]([a-z0-9]|[.](?![.]))*[a-z0-9]@[a-z0-9.-]+\.[a-z]{2,}$');

ALTER TABLE app.users ADD CONSTRAINT users_full_name_length_check CHECK (char_length(trim(full_name)) BETWEEN 2 AND 50);
ALTER TABLE app.users ADD CONSTRAINT users_full_name_no_leading_trailing_spaces CHECK (full_name = trim(full_name));

ALTER TABLE app.users ADD CONSTRAINT users_user_workout_score_check CHECK (user_workout_score >= 0);

CREATE UNIQUE INDEX idx_users_unique_lower_mail ON app.users (LOWER(mail));

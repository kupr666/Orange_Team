CREATE TABLE IF NOT EXISTS app.users (
    id          SERIAL PRIMARY KEY,
    mail        VARCHAR(30)  NOT NULL,
    pass_hash   VARCHAR(255) NOT NULL,
    full_name   VARCHAR(50)  NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL,
    updated_at  TIMESTAMPTZ  NOT NULL
);

ALTER TABLE app.users ADD CONSTRAINT mail_formatcheck 
    CHECK (mail ~* '^[A-Z0-9.%+-]+@[A-Z0-9.-]+.[A-Z]{2,}$');

ALTER TABLE app.users ADD CONSTRAINT full_name_length_check 
    CHECK (char_length(trim(full_name)) BETWEEN 2 AND 50);

ALTER TABLE app.users ADD CONSTRAINT full_name_no_leading_trailing_spaces 
    CHECK (full_name = trim(full_name));

CREATE UNIQUE INDEX idx_users_unique_lower_mail ON app.users (LOWER(mail));

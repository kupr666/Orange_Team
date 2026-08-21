-- app.users intentionally remains unchanged. For website authentication I would:
--   * rename mail to email and increase its limit to 254 characters;
--   * store email in a normalized form before applying the unique constraint;
--   * rename full_name to name to match the domain terminology;
--   * use NOT NULL DEFAULT NOW() for both created_at and updated_at.
CREATE TABLE IF NOT EXISTS app.users (
    id          SERIAL PRIMARY KEY,
    mail        VARCHAR(30)  NOT NULL UNIQUE,
    pass_hash   VARCHAR(255) NOT NULL,
    full_name   VARCHAR(50)  NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

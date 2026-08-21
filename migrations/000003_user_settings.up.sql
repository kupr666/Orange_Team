CREATE TABLE app.user_settings (
    user_id      INTEGER     PRIMARY KEY REFERENCES app.users (id) ON DELETE CASCADE,
    timezone     VARCHAR(64) NOT NULL,
    language     VARCHAR(16) NOT NULL,
    weight_unit  VARCHAR(8)  NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL,
    updated_at   TIMESTAMPTZ NOT NULL
);

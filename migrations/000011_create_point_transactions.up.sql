CREATE TABLE app.point_transactions (
    id           SERIAL      PRIMARY KEY,
    user_id      INTEGER     NOT NULL REFERENCES app.users (id) ON DELETE CASCADE,
    source_type  VARCHAR(64) NOT NULL,
    source_id    INTEGER,
    amount       INTEGER     NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX point_transactions_user_created_at_idx
    ON app.point_transactions (user_id, created_at DESC);

CREATE INDEX point_transactions_source_idx
    ON app.point_transactions (source_type, source_id);

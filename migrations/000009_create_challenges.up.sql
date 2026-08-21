CREATE TABLE app.challenges (
    id             SERIAL       PRIMARY KEY,
    title          VARCHAR(150) NOT NULL,
    description    TEXT,
    period_type    VARCHAR(16)  NOT NULL,
    target_type    VARCHAR(64)  NOT NULL,
    target_value   INTEGER      NOT NULL CHECK (target_value > 0),
    reward_points  INTEGER      NOT NULL DEFAULT 0 CHECK (reward_points >= 0),
    is_active      BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT challenges_period_type_check
        CHECK (period_type IN ('daily', 'weekly'))
);

CREATE INDEX challenges_active_period_type_idx
    ON app.challenges (is_active, period_type);

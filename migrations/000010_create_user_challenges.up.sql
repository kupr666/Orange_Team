CREATE TABLE app.user_challenges (
    id             SERIAL      PRIMARY KEY,
    user_id        INTEGER     NOT NULL REFERENCES app.users (id) ON DELETE CASCADE,
    challenge_id   INTEGER     NOT NULL REFERENCES app.challenges (id) ON DELETE RESTRICT,
    progress       INTEGER     NOT NULL DEFAULT 0 CHECK (progress >= 0),
    status         VARCHAR(16) NOT NULL DEFAULT 'active',
    assigned_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at     TIMESTAMPTZ NOT NULL,
    completed_at   TIMESTAMPTZ,

    CONSTRAINT user_challenges_status_check
        CHECK (status IN ('active', 'completed', 'expired')),
    CONSTRAINT user_challenges_completion_status_check
        CHECK (
            (status = 'completed' AND completed_at IS NOT NULL)
            OR (status IN ('active', 'expired') AND completed_at IS NULL)
        ),
    CONSTRAINT user_challenges_expiration_check
        CHECK (expires_at > assigned_at),
    CONSTRAINT user_challenges_completion_check
        CHECK (completed_at IS NULL OR completed_at >= assigned_at)
);

CREATE INDEX user_challenges_user_status_idx
    ON app.user_challenges (user_id, status);

CREATE INDEX user_challenges_challenge_id_idx
    ON app.user_challenges (challenge_id);

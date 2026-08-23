BEGIN;

-- Temporary development user matching the uuid.Nil placeholder currently used
-- by GET /workouts. Remove this seed once authentication supplies a real user ID.
INSERT INTO app.users (
    id,
    mail,
    pass_hash,
    full_name,
    created_at,
    sex,
    weight_grams,
    birth_date,
    height_cm
)
VALUES (
    '00000000-0000-0000-0000-000000000000',
    'workouts@example.com',
    'development-only-not-a-real-password-hash',
    'Workout Tester',
    NOW(),
    'unspecified',
    75000,
    DATE '2000-01-01',
    175
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO app.workouts (
    id,
    user_id,
    status,
    started_at,
    completed_at,
    created_at,
    updated_at,
    workout_score,
    intensity,
    personal_score_coefficient
)
VALUES
    (
        '10000000-0000-0000-0000-000000000001',
        '00000000-0000-0000-0000-000000000000',
        'planned',
        NULL,
        NULL,
        NOW() - INTERVAL '3 days',
        NOW() - INTERVAL '3 days',
        0,
        NULL,
        5
    ),
    (
        '10000000-0000-0000-0000-000000000002',
        '00000000-0000-0000-0000-000000000000',
        'in_progress',
        NOW() - INTERVAL '45 minutes',
        NULL,
        NOW() - INTERVAL '2 days',
        NOW() - INTERVAL '45 minutes',
        0,
        6,
        5
    ),
    (
        '10000000-0000-0000-0000-000000000003',
        '00000000-0000-0000-0000-000000000000',
        'completed',
        NOW() - INTERVAL '1 day 1 hour',
        NOW() - INTERVAL '1 day',
        NOW() - INTERVAL '1 day 2 hours',
        NOW() - INTERVAL '1 day',
        850,
        8,
        6
    ),
    (
        '10000000-0000-0000-0000-000000000004',
        '00000000-0000-0000-0000-000000000000',
        'cancelled',
        NULL,
        NULL,
        NOW() - INTERVAL '4 days',
        NOW() - INTERVAL '2 days',
        0,
        NULL,
        4
    )
ON CONFLICT (id) DO NOTHING;

COMMIT;

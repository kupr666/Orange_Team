BEGIN;

-- Очищаем таблицу тренировок
TRUNCATE app.workouts CASCADE;

-- Тестовый пользователь
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

-- Вставляем 8 тренировок для различных сценариев
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
    -- 1. planned → in_progress (успех)
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
    -- 2. planned → cancelled (успех)
    (
        '10000000-0000-0000-0000-000000000002',
        '00000000-0000-0000-0000-000000000000',
        'planned',
        NULL,
        NULL,
        NOW() - INTERVAL '2 days',
        NOW() - INTERVAL '2 days',
        0,
        NULL,
        6
    ),
    -- 3. in_progress → completed с intensity (успех)
    (
        '10000000-0000-0000-0000-000000000003',
        '00000000-0000-0000-0000-000000000000',
        'in_progress',
        NOW() - INTERVAL '1 hour',
        NULL,
        NOW() - INTERVAL '1 day',
        NOW() - INTERVAL '1 hour',
        0,
        6,
        7
    ),
    -- 4. in_progress → completed без intensity (ошибка)
    (
        '10000000-0000-0000-0000-000000000004',
        '00000000-0000-0000-0000-000000000000',
        'in_progress',
        NOW() - INTERVAL '45 minutes',
        NULL,
        NOW() - INTERVAL '12 hours',
        NOW() - INTERVAL '45 minutes',
        0,
        NULL,
        5
    ),
    -- 5. in_progress → попытка установить intensity (ошибка)
    (
        '10000000-0000-0000-0000-000000000005',
        '00000000-0000-0000-0000-000000000000',
        'in_progress',
        NOW() - INTERVAL '20 minutes',
        NULL,
        NOW() - INTERVAL '6 hours',
        NOW() - INTERVAL '20 minutes',
        0,
        NULL,
        4
    ),
    -- 6. completed → изменение intensity (успех) и запрет смены статуса (ошибка)
    (
        '10000000-0000-0000-0000-000000000006',
        '00000000-0000-0000-0000-000000000000',
        'completed',
        NOW() - INTERVAL '3 hours',
        NOW() - INTERVAL '2 hours',
        NOW() - INTERVAL '1 day 3 hours',
        NOW() - INTERVAL '2 hours',
        150,
        8,
        6
    ),
    -- 7. cancelled → запрет любых изменений (ошибка)
    (
        '10000000-0000-0000-0000-000000000007',
        '00000000-0000-0000-0000-000000000000',
        'cancelled',
        NULL,
        NULL,
        NOW() - INTERVAL '5 days',
        NOW() - INTERVAL '3 days',
        0,
        NULL,
        4
    ),
    -- 8. in_progress → изменение started_at (успех)
    (
        '10000000-0000-0000-0000-000000000008',
        '00000000-0000-0000-0000-000000000000',
        'in_progress',
        NOW() - INTERVAL '10 minutes',
        NULL,
        NOW() - INTERVAL '1 hour',
        NOW() - INTERVAL '10 minutes',
        0,
        5,
        5
    )
ON CONFLICT (id) DO NOTHING;

COMMIT;
ALTER TABLE app.users
    ALTER COLUMN id SET DEFAULT gen_random_uuid(),
    ALTER COLUMN created_at SET DEFAULT NOW(),
    ALTER COLUMN sex SET DEFAULT 'unspecified',
    ALTER COLUMN weight_grams DROP NOT NULL,
    ALTER COLUMN birth_date DROP NOT NULL,
    ALTER COLUMN height_cm DROP NOT NULL;

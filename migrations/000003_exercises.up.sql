CREATE TABLE IF NOT EXISTS app.exercises (
    id           UUID          PRIMARY KEY,
    version      BIGINT        NOT NULL DEFAULT 1,
    name         VARCHAR(100)  NOT NULL,
    description  VARCHAR(1000) NOT NULL,
    difficulty   SMALLINT      NOT NULL,
    created_at   TIMESTAMPTZ   NOT NULL,
    updated_at   TIMESTAMPTZ,
    type         VARCHAR(20)   NOT NULL
);

ALTER TABLE app.exercises ADD CONSTRAINT exercises_name_unique UNIQUE (name);
ALTER TABLE app.exercises ADD CONSTRAINT exercises_name_no_leading_trailing_spaces CHECK (name = trim(name));
ALTER TABLE app.exercises ADD CONSTRAINT exercises_name_characters CHECK (name ~* '^[a-zA-Zа-яА-ЯёЁ0-9 .,\-()/+="'']+$');
ALTER TABLE app.exercises ADD CONSTRAINT exercises_name_length CHECK (char_length(name) BETWEEN 3 AND 100);

ALTER TABLE app.exercises ADD CONSTRAINT exercises_description_no_leading_trailing_spaces CHECK (description = trim(description));
ALTER TABLE app.exercises ADD CONSTRAINT exercises_description_characters CHECK (description ~* '^[a-zA-Zа-яА-ЯёЁ0-9 .,\-()/+="'']+$');
ALTER TABLE app.exercises ADD CONSTRAINT exercises_description_length CHECK (char_length(description) BETWEEN 1 AND 1000);

ALTER TABLE app.exercises ADD CONSTRAINT exercises_difficulty_check CHECK (difficulty BETWEEN 1 AND 10);

ALTER TABLE app.exercises ADD CONSTRAINT exercises_type_check CHECK (type IN ('weight', 'duration'));

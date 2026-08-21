CREATE TABLE IF NOT EXISTS app.exercises (
    id           SERIAL        PRIMARY KEY,
    name         VARCHAR(100)  NOT NULL,
    description  VARCHAR(1000) NOT NULL,
    base_points  INTEGER       NOT NULL,
    created_at   TIMESTAMPTZ   NOT NULL,
    updated_at   TIMESTAMPTZ
);

ALTER TABLE app.exercises ADD CONSTRAINT exercises_name_unique UNIQUE (name);
ALTER TABLE app.exercises ADD CONSTRAINT exercises_name_no_leading_trailing_spaces CHECK (name = trim(name));
ALTER TABLE app.exercises ADD CONSTRAINT exercises_name_characters CHECK (name ~* '^[A-Za-zА-Яа-яёЁ ]+$');
ALTER TABLE app.exercises ADD CONSTRAINT exercises_name_length CHECK (char_length(name) BETWEEN 3 AND 30);

ALTER TABLE app.exercises ADD CONSTRAINT exercises_description_no_leading_trailing_spaces CHECK (description = trim(description));
ALTER TABLE app.exercises ADD CONSTRAINT exercises_description_characters CHECK (description ~* '^[A-Za-zА-Яа-яёЁ ]+$');
ALTER TABLE app.exercises ADD CONSTRAINT exercises_description_length CHECK (char_length(description) BETWEEN 1 AND 1000);

ALTER TABLE app.exercises ADD CONSTRAINT exercises_base_points_check CHECK (base_points BETWEEN 1 AND 100);

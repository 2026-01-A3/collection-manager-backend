CREATE TABLE IF NOT EXISTS binary_objects
(
    id        INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    base64    TEXT NOT NULL,
    filename  TEXT NOT NULL,
    extension TEXT NOT NULL
);

ALTER TABLE collections
    ADD COLUMN IF NOT EXISTS binary_object_id INTEGER REFERENCES binary_objects (id) ON DELETE SET NULL;

ALTER TABLE collections
    DROP COLUMN IF EXISTS image_url;

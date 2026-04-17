ALTER TABLE collections
    ADD COLUMN IF NOT EXISTS image_url TEXT;

ALTER TABLE collections
    DROP COLUMN IF EXISTS binary_object_id;

DROP TABLE IF EXISTS binary_objects;

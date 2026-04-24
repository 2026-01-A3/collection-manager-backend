CREATE TABLE IF NOT EXISTS items
(
    id               INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name             TEXT           NOT NULL,
    price            NUMERIC(10, 2) NOT NULL DEFAULT 0,
    collection_id    INTEGER        NOT NULL REFERENCES collections (id) ON DELETE CASCADE,
    binary_object_id INTEGER REFERENCES binary_objects (id) ON DELETE SET NULL
);

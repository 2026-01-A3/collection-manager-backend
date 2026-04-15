CREATE TABLE IF NOT EXISTS collections
(
    id
    INTEGER
    GENERATED
    ALWAYS AS
    IDENTITY
    PRIMARY
    KEY,
    name
    TEXT
    NOT
    NULL,
    category_id
    INTEGER
    NOT
    NULL
    REFERENCES
    categories(id),
    image_url
    TEXT
);

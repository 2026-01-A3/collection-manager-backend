CREATE TABLE IF NOT EXISTS categories
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
    NULL
);
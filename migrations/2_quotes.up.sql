CREATE TABLE IF NOT EXISTS quote (
    id SERIAL PRIMARY KEY,
    comment VARCHAR NOT NULL,
    authorId BIGINT
)
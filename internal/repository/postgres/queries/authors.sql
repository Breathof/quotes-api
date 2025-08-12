-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $1 LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: CreateAuthor :one
INSERT INTO authors (
    name, bio
) VALUES (
    $1, $2
)
RETURNING *;

-- name: UpdateAuthor :one
UPDATE authors
SET 
    name = $2,
    bio = $3
WHERE id = $1
RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1;

-- name: CountAuthors :one
SELECT COUNT(*) FROM authors;

-- name: SearchAuthorsByName :many
SELECT * FROM authors
WHERE name ILIKE '%' || $1 || '%'
ORDER BY name
LIMIT $2 OFFSET $3;

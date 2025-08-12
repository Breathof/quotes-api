-- name: GetQuote :one
SELECT 
    q.*,
    a.id as author_id,
    a.name as author_name,
    a.bio as author_bio,
    a.created_at as author_created_at,
    a.updated_at as author_updated_at
FROM quotes q
JOIN authors a ON q.author_id = a.id
WHERE q.id = $1 LIMIT 1;

-- name: ListQuotes :many
SELECT 
    q.*,
    a.id as author_id,
    a.name as author_name,
    a.bio as author_bio,
    a.created_at as author_created_at,
    a.updated_at as author_updated_at
FROM quotes q
JOIN authors a ON q.author_id = a.id
ORDER BY q.created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListQuotesByAuthor :many
SELECT 
    q.*,
    a.id as author_id,
    a.name as author_name,
    a.bio as author_bio,
    a.created_at as author_created_at,
    a.updated_at as author_updated_at
FROM quotes q
JOIN authors a ON q.author_id = a.id
WHERE q.author_id = $1
ORDER BY q.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateQuote :one
INSERT INTO quotes (
    content, author_id, source, tags
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateQuote :one
UPDATE quotes
SET 
    content = $2,
    author_id = $3,
    source = $4,
    tags = $5
WHERE id = $1
RETURNING *;

-- name: DeleteQuote :exec
DELETE FROM quotes
WHERE id = $1;

-- name: CountQuotes :one
SELECT COUNT(*) FROM quotes;

-- name: SearchQuotesByContent :many
SELECT 
    q.*,
    a.id as author_id,
    a.name as author_name,
    a.bio as author_bio,
    a.created_at as author_created_at,
    a.updated_at as author_updated_at
FROM quotes q
JOIN authors a ON q.author_id = a.id
WHERE q.content ILIKE '%' || $1 || '%'
ORDER BY q.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetRandomQuote :one
SELECT 
    q.*,
    a.id as author_id,
    a.name as author_name,
    a.bio as author_bio,
    a.created_at as author_created_at,
    a.updated_at as author_updated_at
FROM quotes q
JOIN authors a ON q.author_id = a.id
ORDER BY RANDOM()
LIMIT 1;

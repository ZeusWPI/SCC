-- name: MessageGetLast :one
SELECT *
FROM message
ORDER BY id DESC
LIMIT 1;

-- name: MessageGetByID :one
SELECT *
FROM message
WHERE id = $1;

-- name: MessageGetSinceID :many
SELECT *
FROM message
WHERE id > $1
ORDER BY created_at ASC;

-- name: MessageCreate :one
INSERT INTO message (name, ip, message)
VALUES ($1, $2, $3)
RETURNING id;

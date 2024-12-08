-- CRUD

-- name: GetAllMessages :many
SELECT *
FROM message;

-- name: GetMessageByID :one
SELECT *
FROM message
WHERE id = $1;

-- name: CreateMessage :one
INSERT INTO message (name, ip, message)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateMessage :one
UPDATE message
SET name = $1, ip = $2, message = $3
WHERE id = $4
RETURNING *;

-- name: DeleteMessage :execrows
DELETE FROM message
WHERE id = $1;


-- Other


-- name: GetLastMessage :one
SELECT *
FROM message
ORDER BY id DESC
LIMIT 1;

-- name: GetMessageSinceID :many
SELECT *
FROM message
WHERE id > $1
ORDER BY created_at ASC;

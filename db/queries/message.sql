-- CRUD

-- name: GetAllMessages :many
SELECT *
FROM message;

-- name: GetMessageByID :one
SELECT *
FROM message
WHERE id = ?;

-- name: CreateMessage :one
INSERT INTO message (name, ip, message)
VALUES (?, ?, ?)
RETURNING *;

-- name: UpdateMessage :one
UPDATE message
SET name = ?, ip = ?, message = ?
WHERE id = ?
RETURNING *;

-- name: DeleteMessage :execrows
DELETE FROM message
WHERE id = ?;


-- Other


-- name: GetLastMessage :one
SELECT *
FROM message
ORDER BY id DESC
LIMIT 1;

-- name: GetMessageSinceID :many
SELECT *
FROM message
WHERE id > ?
ORDER BY created_at ASC;

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

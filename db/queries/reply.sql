-- name: ReplyGetSinceMessageID :many
SELECT *
FROM reply
WHERE message_id >= $1
ORDER BY message_id;

-- name: ReplyCreate :one
INSERT INTO reply (message_id, name, message)
VALUES ($1, $2, $3)
RETURNING id;

-- CRUD

-- name: CreateSongHistory :one
INSERT INTO song_history (song_id)
VALUES (?)
RETURNING *;


-- Other


-- name: GetLastSongHistory :one
SELECT *
FROM song_history
ORDER BY created_at DESC
LIMIT 1;

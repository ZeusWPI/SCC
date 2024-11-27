-- CRUD

-- name: GetAllSongs :many
SELECT *
FROM song;

-- name: GetSongByID :one
SELECT *
FROM song
WHERE id = ?;

-- name: CreateSong :one
INSERT INTO song (title, artists, spotify_id, duration_ms)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: UpdateSong :one
UPDATE song
SET title = ?, artists = ?, spotify_id = ?, duration_ms = ?
WHERE id = ?
RETURNING *;

-- name: DeleteSong :execrows
DELETE FROM song
WHERE id = ?;


-- Other

-- name: GetSongBySpotifyID :one
SELECT *
FROM song
WHERE spotify_id = ?;

-- CRUD

-- name: GetAllSpotify :many
SELECT *
FROM spotify;

-- name: GetSpotifyByID :one
SELECT *
FROM spotify
WHERE id = ?;

-- name: CreateSpotify :one
INSERT INTO spotify (title, artists, spotify_id, duration_ms)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: UpdateSpotify :one
UPDATE spotify
SET title = ?, artists = ?, spotify_id = ?, duration_ms = ?
WHERE id = ?
RETURNING *;

-- name: DeleteSpotify :execrows
DELETE FROM spotify
WHERE id = ?;


-- Other

-- name: GetSpotifyBySpotifyID :one
SELECT *
FROM spotify
WHERE spotify_id = ?;

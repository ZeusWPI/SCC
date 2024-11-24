-- CRUD

-- name: GetAllSeasons :many
SELECT *
FROM season;

-- name: GetSeasonByID :one
SELECT *
FROM season
WHERE id = ?;

-- name: CreateSeason :one
INSERT INTO season (name, start, end, current)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: UpdateSeason :one
UPDATE season
SET name = ?, start = ?, end = ?, current = ?
WHERE id = ?
RETURNING *;

-- name: DeleteSeason :execrows
DELETE FROM season
WHERE id = ?;

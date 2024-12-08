-- CRUD

-- name: GetAllSeasons :many
SELECT *
FROM season;

-- name: GetSeasonByID :one
SELECT *
FROM season
WHERE id = $1;

-- name: CreateSeason :one
INSERT INTO season (name, start, "end", current)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateSeason :one
UPDATE season
SET name = $1, start = $2, "end" = $3, current = $4
WHERE id = $5
RETURNING *;

-- name: DeleteSeason :execrows
DELETE FROM season
WHERE id = $1;

-- name: DeleteSeasonAll :execrows
DELETE FROM season;



-- Other


-- name: GetSeasonCurrent :one
SELECT *
FROM season
WHERE current = true;

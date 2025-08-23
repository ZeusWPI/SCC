-- name: SeasonGetAll :many
SELECT *
FROM season;

-- name: SeasonGetCurrent :one
SELECT *
FROM season
WHERE current = true;

-- name: SeasonCreate :one
INSERT INTO season (name, start, "end", current)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: SeasonUpdate :exec
UPDATE season
SET name = $1, start = $2, "end" = $3, current = $4
WHERE id = $5
RETURNING *;

-- name: SeasonDeleteAll :exec
DELETE FROM season;

-- CRUD

-- name: GetAllScans :many
SELECT *
FROM scan;

-- name: GetScanByID :one
SELECT *
FROM scan
WHERE id = $1;

-- name: CreateScan :one
INSERT INTO scan (scan_id, scan_time)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateScan :one
UPDATE scan
SET scan_id = $1, scan_time = $2
WHERE id = $3
RETURNING *;

-- name: DeleteScan :execrows
DELETE FROM scan
WHERE id = $1;


-- Other


-- name: GetLastScan :one
SELECT *
FROM scan
ORDER BY id DESC
LIMIT 1;

-- name: GetAllScansSinceID :many
SELECT *
FROM scan
WHERE id > $1
ORDER BY scan_id, scan_time ASC;

-- name: GetScansInCurrentSeason :one
SELECT COUNT(*) AS amount
FROM scan
WHERE scan_time >= (SELECT start_date FROM season WHERE current = true) AND
        scan_time <= (SELECT end_date FROM season WHERE current = true);

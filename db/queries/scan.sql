-- CRUD

-- name: GetAllScans :many
SELECT *
FROM scan;

-- name: GetScanByID :one
SELECT *
FROM scan
WHERE id = ?;

-- name: CreateScan :one
INSERT INTO scan (scan_time)
VALUES (?)
RETURNING *;

-- name: UpdateScan :one
UPDATE scan
SET scan_time = ?
WHERE id = ?
RETURNING *;

-- name: DeleteScan :execrows
DELETE FROM scan
WHERE id = ?;


-- Other


-- name: GetLastScan :one
SELECT *
FROM scan
ORDER BY id DESC
LIMIT 1;

-- name: GetAllScansSinceID :many
SELECT *
FROM scan
WHERE id > ?
ORDER BY scan_time ASC;

-- name: GetScansInCurrentSeason :one
SELECT COUNT(*) AS amount
FROM scan
WHERE scan_time >= (SELECT start_date FROM season WHERE current = true) AND
        scan_time <= (SELECT end_date FROM season WHERE current = true);

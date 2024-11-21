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


-- name: GetLatestScan :one
SELECT *
FROM scan
ORDER BY id DESC
LIMIT 1;

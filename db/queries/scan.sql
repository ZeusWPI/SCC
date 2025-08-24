-- name: ScanGetLast :one
SELECT *
FROM scan
ORDER BY id DESC
LIMIT 1;

-- name: ScanGetAllSinceID :many
SELECT *
FROM scan
WHERE id > $1
ORDER BY scan_id, scan_time ASC;

-- name: ScanGetInSeason :many
SELECT sc.*
FROM scan sc
LEFT JOIN season se ON se.start <= sc.scan_time AND sc.scan_time < se.end
WHERE se.id = $1
ORDER BY sc.scan_time ASC;

-- name: ScanCreate :one
INSERT INTO scan (scan_id, scan_time)
VALUES ($1, $2)
RETURNING id;


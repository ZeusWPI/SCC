-- CRUD

-- name: GetAllTaps :many
SELECT *
FROM tap;

-- name: GetTapByID :one
SELECT *
FROM tap
WHERE id = ?;

-- name: CreateTap :one
INSERT INTO tap (order_id, order_created_at, name, category)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: UpdateTap :one
UPDATE tap
SET order_id = ?, order_created_at = ?, name = ?, category = ?
WHERE id = ?
RETURNING *;

-- name: DeleteTap :execrows
DELETE FROM tap
WHERE id = ?;


-- Other


-- name: GetTapByOrderID :one
SELECT *
FROM tap
WHERE order_id = ?;

-- name: GetTapByCategory :many
SELECT *
FROM tap
WHERE category = ?;

-- name: GetLastOrderByOrderID :one
SELECT *
FROM tap
ORDER BY order_id DESC
LIMIT 1;

-- name: GetOrderCount :many
SELECT category, COUNT(*)
FROM tap
GROUP BY category;

-- name: GetOrderCountByCategorySinceOrderID :many
SELECT category, COUNT(*)
FROM tap
WHERE order_id >= ?
GROUP BY category;

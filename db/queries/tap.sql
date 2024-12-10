-- CRUD

-- name: GetAllTaps :many
SELECT *
FROM tap;

-- name: GetTapByID :one
SELECT *
FROM tap
WHERE id = $1;

-- name: CreateTap :one
INSERT INTO tap (order_id, order_created_at, name, category)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateTap :one
UPDATE tap
SET order_id = $1, order_created_at = $2, name = $3, category = $4
WHERE id = $5
RETURNING *;

-- name: DeleteTap :execrows
DELETE FROM tap
WHERE id = $1;


-- Other


-- name: GetTapByOrderID :one
SELECT *
FROM tap
WHERE order_id = $1;

-- name: GetTapByCategory :many
SELECT *
FROM tap
WHERE category = $1;

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
SELECT category, COUNT(*), MAX(order_created_at)::TIMESTAMP AS latest_order_created_at
FROM tap
WHERE order_id >= $1
GROUP BY category;

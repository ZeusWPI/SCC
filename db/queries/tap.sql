-- name: TapGetLast :one
SELECT *
FROM tap
ORDER BY order_id DESC
LIMIT 1;

-- name: TapGetCountByCategory :many
SELECT category, COUNT(*), MAX(order_created_at)::TIMESTAMP AS latest_order_created_at
FROM tap
GROUP BY category;

-- name: TapCreate :one
INSERT INTO tap (order_id, order_created_at, name, category)
VALUES ($1, $2, $3, $4)
RETURNING id;


-- name: GamificationGetAll :many
SELECT *
FROM gamification
ORDER BY score DESC;

-- name: GamificationCreate :one
INSERT INTO gamification (name, score, avatar)
VALUES ($1, $2, $3)
RETURNING id;

-- name: GamificationDeleteAll :exec
DELETE FROM gamification;


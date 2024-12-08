-- CRUD

-- name: GetAllGamification :many
SELECT *
FROM gamification;

-- name: CreateGamification :one
INSERT INTO gamification (name, score, avatar)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteGamification :execrows
DELETE FROM gamification
WHERE id = $1;

-- name: DeleteGamificationAll :execrows
DELETE FROM gamification;


-- Other


-- name: UpdateGamificationScore :one
UPDATE gamification
SET score = $1
WHERE id = $2
RETURNING *;

-- name: GetAllGamificationByScore :many
SELECT *
FROM gamification
ORDER BY score DESC;

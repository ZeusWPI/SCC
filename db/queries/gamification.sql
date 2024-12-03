-- CRUD

-- name: GetAllGamification :many
SELECT *
FROM gamification;

-- name: CreateGamification :one
INSERT INTO gamification (name, score, avatar)
VALUES (?, ?, ?)
RETURNING *;

-- name: DeleteGamification :execrows
DELETE FROM gamification
WHERE id = ?;

-- name: DeleteGamificationAll :execrows
DELETE FROM gamification;


-- Other


-- name: UpdateGamificationScore :one
UPDATE gamification
SET score = ?
WHERE id = ?
RETURNING *;

-- name: GetAllGamificationByScore :many
SELECT *
FROM gamification
ORDER BY score DESC;

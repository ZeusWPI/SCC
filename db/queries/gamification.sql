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


-- Other


-- name: UpdateGamificationScore :one
UPDATE gamification
SET score = ?
WHERE id = ?
RETURNING *;

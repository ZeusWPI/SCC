-- CRUD


-- name: GetAllEvents :many
SELECT *
FROM event;

-- name: CreateEvent :one
INSERT INTO event (name, date, academic_year, location)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: DeleteEvent :exec
DELETE FROM event
WHERE id = ?;


-- Other


-- name: GetEventByAcademicYear :many
SELECT *
FROM event
WHERE academic_year = ?;

-- name: DeleteEventByAcademicYear :exec
DELETE FROM event
WHERE academic_year = ?;

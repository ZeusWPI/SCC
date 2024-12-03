-- CRUD


-- name: GetAllEvents :many
SELECT *
FROM event;

-- name: CreateEvent :one
INSERT INTO event (name, date, academic_year, location, poster)
VALUES (?, ?, ?, ?, ?)
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

-- name: GetEventsCurrentAcademicYear :many
SELECT *
FROM event
WHERE academic_year = (
    SELECT MAX(academic_year)
    FROM event
)
ORDER BY date ASC;

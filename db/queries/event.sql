-- CRUD


-- name: GetAllEvents :many
SELECT *
FROM event;

-- name: CreateEvent :one
INSERT INTO event (name, date, academic_year, location, poster)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteEvent :exec
DELETE FROM event
WHERE id = $1;


-- Other


-- name: GetEventByAcademicYear :many
SELECT *
FROM event
WHERE academic_year = $1;

-- name: DeleteEventByAcademicYear :exec
DELETE FROM event
WHERE academic_year = $1;

-- name: GetEventsCurrentAcademicYear :many
SELECT *
FROM event
WHERE academic_year = (
    SELECT MAX(academic_year)
    FROM event
)
ORDER BY date ASC;

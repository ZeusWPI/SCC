-- name: EventGetAllByYear :many
SELECT *
FROM event
WHERE academic_year = $1;

-- name: EventGetAllByCurrentYear :many
SELECT *
FROM event
WHERE academic_year = (
    SELECT MAX(academic_year)
    FROM event
)
ORDER BY date ASC;

-- name: EventCreate :one
INSERT INTO event (name, date, academic_year, location, poster)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;


-- name: EventDeleteAll :exec
DELETE FROM event;

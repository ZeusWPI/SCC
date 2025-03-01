// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: season.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createSeason = `-- name: CreateSeason :one
INSERT INTO season (name, start, "end", current)
VALUES ($1, $2, $3, $4)
RETURNING id, name, start, "end", current
`

type CreateSeasonParams struct {
	Name    string
	Start   pgtype.Timestamp
	End     pgtype.Timestamp
	Current bool
}

func (q *Queries) CreateSeason(ctx context.Context, arg CreateSeasonParams) (Season, error) {
	row := q.db.QueryRow(ctx, createSeason,
		arg.Name,
		arg.Start,
		arg.End,
		arg.Current,
	)
	var i Season
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Start,
		&i.End,
		&i.Current,
	)
	return i, err
}

const deleteSeason = `-- name: DeleteSeason :execrows
DELETE FROM season
WHERE id = $1
`

func (q *Queries) DeleteSeason(ctx context.Context, id int32) (int64, error) {
	result, err := q.db.Exec(ctx, deleteSeason, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const deleteSeasonAll = `-- name: DeleteSeasonAll :execrows
DELETE FROM season
`

func (q *Queries) DeleteSeasonAll(ctx context.Context) (int64, error) {
	result, err := q.db.Exec(ctx, deleteSeasonAll)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getAllSeasons = `-- name: GetAllSeasons :many

SELECT id, name, start, "end", current
FROM season
`

// CRUD
func (q *Queries) GetAllSeasons(ctx context.Context) ([]Season, error) {
	rows, err := q.db.Query(ctx, getAllSeasons)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Season
	for rows.Next() {
		var i Season
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Start,
			&i.End,
			&i.Current,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSeasonByID = `-- name: GetSeasonByID :one
SELECT id, name, start, "end", current
FROM season
WHERE id = $1
`

func (q *Queries) GetSeasonByID(ctx context.Context, id int32) (Season, error) {
	row := q.db.QueryRow(ctx, getSeasonByID, id)
	var i Season
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Start,
		&i.End,
		&i.Current,
	)
	return i, err
}

const getSeasonCurrent = `-- name: GetSeasonCurrent :one


SELECT id, name, start, "end", current
FROM season
WHERE current = true
`

// Other
func (q *Queries) GetSeasonCurrent(ctx context.Context) (Season, error) {
	row := q.db.QueryRow(ctx, getSeasonCurrent)
	var i Season
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Start,
		&i.End,
		&i.Current,
	)
	return i, err
}

const updateSeason = `-- name: UpdateSeason :one
UPDATE season
SET name = $1, start = $2, "end" = $3, current = $4
WHERE id = $5
RETURNING id, name, start, "end", current
`

type UpdateSeasonParams struct {
	Name    string
	Start   pgtype.Timestamp
	End     pgtype.Timestamp
	Current bool
	ID      int32
}

func (q *Queries) UpdateSeason(ctx context.Context, arg UpdateSeasonParams) (Season, error) {
	row := q.db.QueryRow(ctx, updateSeason,
		arg.Name,
		arg.Start,
		arg.End,
		arg.Current,
		arg.ID,
	)
	var i Season
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Start,
		&i.End,
		&i.Current,
	)
	return i, err
}

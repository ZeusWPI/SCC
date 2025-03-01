// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: scan.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createScan = `-- name: CreateScan :one
INSERT INTO scan (scan_id, scan_time)
VALUES ($1, $2)
RETURNING id, scan_time, scan_id
`

type CreateScanParams struct {
	ScanID   int32
	ScanTime pgtype.Timestamptz
}

func (q *Queries) CreateScan(ctx context.Context, arg CreateScanParams) (Scan, error) {
	row := q.db.QueryRow(ctx, createScan, arg.ScanID, arg.ScanTime)
	var i Scan
	err := row.Scan(&i.ID, &i.ScanTime, &i.ScanID)
	return i, err
}

const deleteScan = `-- name: DeleteScan :execrows
DELETE FROM scan
WHERE id = $1
`

func (q *Queries) DeleteScan(ctx context.Context, id int32) (int64, error) {
	result, err := q.db.Exec(ctx, deleteScan, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getAllScans = `-- name: GetAllScans :many

SELECT id, scan_time, scan_id
FROM scan
`

// CRUD
func (q *Queries) GetAllScans(ctx context.Context) ([]Scan, error) {
	rows, err := q.db.Query(ctx, getAllScans)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Scan
	for rows.Next() {
		var i Scan
		if err := rows.Scan(&i.ID, &i.ScanTime, &i.ScanID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllScansSinceID = `-- name: GetAllScansSinceID :many
SELECT id, scan_time, scan_id
FROM scan
WHERE id > $1
ORDER BY scan_id, scan_time ASC
`

func (q *Queries) GetAllScansSinceID(ctx context.Context, id int32) ([]Scan, error) {
	rows, err := q.db.Query(ctx, getAllScansSinceID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Scan
	for rows.Next() {
		var i Scan
		if err := rows.Scan(&i.ID, &i.ScanTime, &i.ScanID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLastScan = `-- name: GetLastScan :one


SELECT id, scan_time, scan_id
FROM scan
ORDER BY id DESC
LIMIT 1
`

// Other
func (q *Queries) GetLastScan(ctx context.Context) (Scan, error) {
	row := q.db.QueryRow(ctx, getLastScan)
	var i Scan
	err := row.Scan(&i.ID, &i.ScanTime, &i.ScanID)
	return i, err
}

const getScanByID = `-- name: GetScanByID :one
SELECT id, scan_time, scan_id
FROM scan
WHERE id = $1
`

func (q *Queries) GetScanByID(ctx context.Context, id int32) (Scan, error) {
	row := q.db.QueryRow(ctx, getScanByID, id)
	var i Scan
	err := row.Scan(&i.ID, &i.ScanTime, &i.ScanID)
	return i, err
}

const getScansInCurrentSeason = `-- name: GetScansInCurrentSeason :one
SELECT COUNT(*) AS amount
FROM scan
WHERE scan_time >= (SELECT start_date FROM season WHERE current = true) AND
        scan_time <= (SELECT end_date FROM season WHERE current = true)
`

func (q *Queries) GetScansInCurrentSeason(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, getScansInCurrentSeason)
	var amount int64
	err := row.Scan(&amount)
	return amount, err
}

const updateScan = `-- name: UpdateScan :one
UPDATE scan
SET scan_id = $1, scan_time = $2
WHERE id = $3
RETURNING id, scan_time, scan_id
`

type UpdateScanParams struct {
	ScanID   int32
	ScanTime pgtype.Timestamptz
	ID       int32
}

func (q *Queries) UpdateScan(ctx context.Context, arg UpdateScanParams) (Scan, error) {
	row := q.db.QueryRow(ctx, updateScan, arg.ScanID, arg.ScanTime, arg.ID)
	var i Scan
	err := row.Scan(&i.ID, &i.ScanTime, &i.ScanID)
	return i, err
}

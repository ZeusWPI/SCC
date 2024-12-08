package dto

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zeusWPI/scc/internal/pkg/db/sqlc"
)

// Scan is the DTO for the scan
type Scan struct {
	ID       int32     `json:"id"`
	ScanTime time.Time `json:"scan_time" validate:"required"`
}

// ScanDTO converts a sqlc.Scan to a Scan
func ScanDTO(scan sqlc.Scan) *Scan {
	return &Scan{
		ID:       scan.ID,
		ScanTime: scan.ScanTime.Time,
	}
}

// CreateParams converts a Scan to sqlc.CreateScanParams
func (s *Scan) CreateParams() pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: s.ScanTime}
}

// UpdateParams converts a Scan to sqlc.UpdateScanParams
func (s *Scan) UpdateParams() sqlc.UpdateScanParams {
	return sqlc.UpdateScanParams{
		ID:       s.ID,
		ScanTime: pgtype.Timestamptz{Time: s.ScanTime, Valid: true},
	}
}

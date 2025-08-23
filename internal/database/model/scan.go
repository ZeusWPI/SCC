package model

import (
	"time"

	"github.com/zeusWPI/scc/pkg/sqlc"
)

type Scan struct {
	ID       int
	ScanID   int
	ScanTime time.Time
}

func ScanModel(s sqlc.Scan) *Scan {
	return &Scan{
		ID:       int(s.ID),
		ScanID:   int(s.ScanID),
		ScanTime: s.ScanTime.Time,
	}
}

package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zeusWPI/scc/internal/database/model"
	"github.com/zeusWPI/scc/pkg/sqlc"
	"github.com/zeusWPI/scc/pkg/utils"
)

type Scan struct {
	repo Repository
}

func (r *Repository) NewScan() *Scan {
	return &Scan{
		repo: *r,
	}
}

func (s *Scan) GetLast(ctx context.Context) (*model.Scan, error) {
	scan, err := s.repo.queries(ctx).ScanGetLast(ctx)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get last scan %w", err)
		}
		return nil, nil
	}

	return model.ScanModel(scan), nil
}

func (s *Scan) GetAllSinceID(ctx context.Context, id int) ([]*model.Scan, error) {
	scans, err := s.repo.queries(ctx).ScanGetAllSinceID(ctx, int32(id))
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get scans since id %d | %w", id, err)
		}
		return nil, nil
	}

	return utils.SliceMap(scans, model.ScanModel), nil
}

func (s *Scan) Create(ctx context.Context, scan *model.Scan) error {
	id, err := s.repo.queries(ctx).ScanCreate(ctx, sqlc.ScanCreateParams{
		ScanID:   int32(scan.ScanID),
		ScanTime: pgtype.Timestamptz{Time: scan.ScanTime, Valid: !scan.ScanTime.IsZero()},
	})
	if err != nil {
		return fmt.Errorf("create scan %+v | %w", *scan, err)
	}

	scan.ID = int(id)

	return nil
}

// Package zess provides all zess related logic
package zess

import (
	"context"
	"errors"
	"slices"

	"github.com/jackc/pgx/v5"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/internal/pkg/db/dto"
	"github.com/zeusWPI/scc/internal/pkg/db/sqlc"
	"github.com/zeusWPI/scc/pkg/config"
	"github.com/zeusWPI/scc/pkg/util"
)

// Zess represents a zess instance
type Zess struct {
	db  *db.DB
	api string
}

// New creates a new zess instance
func New(db *db.DB) *Zess {
	return &Zess{
		db:  db,
		api: config.GetDefaultString("backend.zess.api", "https://zess.zeus.gent/api"),
	}
}

// UpdateSeasons updates the seasons
func (z *Zess) UpdateSeasons() error {
	seasons, err := z.db.Queries.GetAllSeasons(context.Background())
	if err != nil && err != pgx.ErrNoRows {
		return err
	}

	// Get all seasons from zess
	zessSeasons, err := z.getSeasons()
	if err != nil {
		return err
	}

	if slices.CompareFunc(util.SliceMap(seasons, dto.SeasonDTO), *zessSeasons, dto.SeasonCmp) == 0 {
		return nil
	}

	// The seasons differ
	// Delete all existing and enter the new ones
	if _, err := z.db.Queries.DeleteSeasonAll(context.Background()); err != nil {
		return err
	}

	var errs []error
	for _, season := range *zessSeasons {
		if _, err := z.db.Queries.CreateSeason(context.Background(), season.CreateParams()); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// UpdateScans updates the scans
func (z *Zess) UpdateScans() error {
	lastScan, err := z.db.Queries.GetLastScan(context.Background())
	if err != nil {
		if err != pgx.ErrNoRows {
			return err
		}

		lastScan = sqlc.Scan{ID: -1}
	}

	// Get all scans
	zessScans, err := z.getScans()
	if err != nil {
		return err
	}

	errs := make([]error, 0)
	for _, scan := range *zessScans {
		if lastScan.ID >= scan.ID {
			continue
		}

		_, err := z.db.Queries.CreateScan(context.Background(), scan.CreateParams())
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

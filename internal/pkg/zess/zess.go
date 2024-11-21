// Package zess provides all zess related logic
package zess

import (
	"context"
	"database/sql"
	"errors"
	"slices"

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
	api := config.GetDefaultString("zess.api", "https://zess.zeus.gent")
	return &Zess{db: db, api: api}
}

// UpdateSeasons updates the seasons
func (z *Zess) UpdateSeasons() error {
	seasons, err := z.db.Queries.GetAllSeasons(context.Background())
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	// Get all seasons from zess
	zessSeasons, err := z.getSeasons()
	if err != nil {
		return err
	}

	equal := slices.CompareFunc(util.SliceMap(seasons, dto.SeasonDTO), *zessSeasons, dto.SeasonCmp)

	// Same seasons
	if equal == 0 {
		return nil
	}

	// Update seasons
	errs := make([]error, 0)

	for i, season := range *zessSeasons {
		if i < len(seasons) {
			// Update seasons
			seasons[i].ID = season.ID
			seasons[i].Name = season.Name
			seasons[i].Start = season.Start
			seasons[i].End = season.End

			_, err := z.db.Queries.UpdateSeason(context.Background(), dto.SeasonDTO(seasons[i]).UpdateParams())
			if err != nil {
				errs = append(errs, err)
			}
		} else {
			// Create seasons
			_, err := z.db.Queries.CreateSeason(context.Background(), season.CreateParams())
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	// Delete seasons
	for i := len(*zessSeasons); i < len(seasons); i++ {
		_, err := z.db.Queries.DeleteSeason(context.Background(), seasons[i].ID)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// UpdateScans updates the scans
func (z *Zess) UpdateScans() error {
	lastScan, err := z.db.Queries.GetLatestScan(context.Background())
	if err != nil {
		if err != sql.ErrNoRows {
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

// Package zess provides all zess related logic
package zess

import (
	"context"
	"slices"

	"github.com/zeusWPI/scc/internal/database/model"
	"github.com/zeusWPI/scc/internal/database/repository"
	"github.com/zeusWPI/scc/pkg/config"
	"github.com/zeusWPI/scc/pkg/utils"
)

type Zess struct {
	repo    repository.Repository
	scan    repository.Scan
	seasons repository.Season
	url     string
}

func New(repo repository.Repository) *Zess {
	return &Zess{
		repo:    repo,
		scan:    *repo.NewScan(),
		seasons: *repo.NewSeason(),
		url:     config.GetDefaultString("backend.zess.url", "https://zess.zeus.gent/api"),
	}
}

func (z *Zess) UpdateSeasons(ctx context.Context) error {
	zessSeasons, err := z.getSeasons(ctx)
	if err != nil {
		return err
	}

	return z.repo.WithRollback(ctx, func(ctx context.Context) error {
		if err := z.seasons.DeleteAll(ctx); err != nil {
			return err
		}

		for _, season := range zessSeasons {
			if err := z.seasons.Create(ctx, &season); err != nil {
				return err
			}
		}

		return nil
	})
}

func (z *Zess) UpdateScans(ctx context.Context) error {
	lastScan, err := z.scan.GetLast(ctx)
	if err != nil {
		return err
	}
	if lastScan == nil {
		lastScan = &model.Scan{ScanID: -1}
	}

	// Get all scans
	zessScans, err := z.getScans(ctx)
	if err != nil {
		return err
	}

	scans := utils.SliceFilter(zessScans, func(s model.Scan) bool { return s.ScanID > lastScan.ScanID })
	slices.SortFunc(scans, func(a, b model.Scan) int { return a.ScanID - b.ScanID })

	for _, scan := range scans {
		if err := z.scan.Create(ctx, &scan); err != nil {
			return err
		}
	}

	return nil
}

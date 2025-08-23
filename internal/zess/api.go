package zess

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zeusWPI/scc/internal/database/model"
	"github.com/zeusWPI/scc/pkg/date"
	"github.com/zeusWPI/scc/pkg/utils"
)

type seasonAPI struct {
	Name    string    `json:"name"`
	Start   date.Date `json:"start"`
	End     date.Date `json:"end"`
	Current bool      `json:"is_current"`
}

func (s seasonAPI) toModel() model.Season {
	return model.Season{
		Name:    s.Name,
		Start:   s.Start,
		End:     s.End,
		Current: s.Current,
	}
}

func (z *Zess) getSeasons(ctx context.Context) ([]model.Season, error) {
	resp, err := utils.DoRequest(ctx, utils.DoRequestValues{
		Method: "GET",
		URL:    fmt.Sprintf("%s/seasons", z.url),
	})
	if err != nil {
		return nil, fmt.Errorf("http get all zess seasons %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	var seasons []seasonAPI
	if err := json.NewDecoder(resp.Body).Decode(&seasons); err != nil {
		return nil, fmt.Errorf("decode http zess seasons %w", err)
	}

	return utils.SliceMap(seasons, func(s seasonAPI) model.Season { return s.toModel() }), nil
}

type scanAPI struct {
	ScanID   int       `json:"scan_id"`
	ScanTime time.Time `json:"scan_time"`
}

func (s scanAPI) toModel() model.Scan {
	return model.Scan{
		ScanID:   s.ScanID,
		ScanTime: s.ScanTime,
	}
}

func (z *Zess) getScans(ctx context.Context) ([]model.Scan, error) {
	resp, err := utils.DoRequest(ctx, utils.DoRequestValues{
		Method: "GET",
		URL:    fmt.Sprintf("%s/recent_scans", z.url),
	})
	if err != nil {
		return nil, fmt.Errorf("http get recent zess scans %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	var scans []scanAPI
	if err := json.NewDecoder(resp.Body).Decode(&scans); err != nil {
		return nil, fmt.Errorf("decode http zess scans %w", err)
	}

	return utils.SliceMap(scans, func(s scanAPI) model.Scan { return s.toModel() }), nil
}

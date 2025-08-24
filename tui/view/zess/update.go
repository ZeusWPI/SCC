package zess

import (
	"context"
	"slices"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/internal/database/model"
	"github.com/zeusWPI/scc/pkg/utils"
	"github.com/zeusWPI/scc/tui/view"
)

func (w week) equal(w2 week) bool {
	return w.scans == w2.scans && w.start.Equal(w2.start)
}

func updateWeeks(ctx context.Context, view view.View) (tea.Msg, error) {
	m := view.(*Model)

	season, err := m.repoSeason.GetCurrent(ctx)
	if err != nil {
		return nil, err
	}
	if season == nil {
		return nil, nil
	}

	scans, err := m.repoScan.GetInSeason(ctx, *season)
	if err != nil {
		return nil, err
	}
	if scans == nil {
		return nil, nil
	}
	slices.SortFunc(scans, func(a, b *model.Scan) int { return a.ScanTime.Compare(b.ScanTime) })

	weekMap := make(map[time.Time]week)

	for _, scan := range scans {
		start := getStartOfWeek(scan.ScanTime)

		entry, ok := weekMap[start]
		if !ok {
			entry = week{
				start: start,
				scans: 0,
			}
		}

		entry.scans++
		weekMap[start] = entry
	}

	weeks := utils.MapValues(weekMap)

	if len(weeks) != len(m.weeks) {
		return Msg{weeks: weeks}, nil
	}

	for idx, week := range weeks {
		if !week.equal(m.weeks[idx]) {
			return Msg{weeks: weeks}, nil
		}
	}

	return nil, nil
}

func getStartOfWeek(t time.Time) time.Time {
	weekDay := int(t.Weekday())

	if weekDay == 0 {
		weekDay = 7
	}

	return time.Date(
		t.Year(), t.Month(), t.Day()-weekDay+1,
		0, 0, 0, 0, t.Location(),
	)
}

// Package zess provides the functions to draw an overview of the zess scans on a TUI
package zess

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/NimbleMarkets/ntcharts/barchart"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/pkg/config"
	"github.com/zeusWPI/scc/ui/view"
	"go.uber.org/zap"
)

// time represents a time object by keeping the year and week number
type time struct {
	year int
	week int
}

type weekScan struct {
	time   time
	amount int64
	label  string
}

// Model represents the Model for the zess view
type Model struct {
	db            *db.DB
	lastScanID    int64
	scans         []weekScan // Queue of scans per week
	maxWeekScans  int64
	currentSeason time // Start week of the season
	seasonScans   int64
}

// Msg is the base message to indicate that something changed in the zess view
type Msg struct{}

// scanMsg is used to indicate that the zess view should be updated with new scans
type scanMsg struct {
	Msg
	lastScanID int64
	scans      []weekScan
}

// seasonMsg is used to indicate that the current season changed.
type seasonMsg struct {
	Msg
	start time
}

// NewModel creates a new zess model view
func NewModel(db *db.DB) view.View {
	z := &Model{
		db:            db,
		lastScanID:    -1,
		scans:         make([]weekScan, 0),
		maxWeekScans:  -1,
		currentSeason: time{year: -1, week: -1},
		seasonScans:   0,
	}

	// Populate with data
	// The order in which this is called is important!
	msgScans, err := updateScans(db, z)
	if err != nil {
		zap.S().Error("TUI: Unable to update zess scans\n", err)
		return z
	}
	_, _ = z.Update(msgScans)

	msgSeason, err := updateSeason(db, z)
	if err != nil {
		zap.S().Error("TUI: Unable to update zess seasons\n", err)
		return z
	}
	_, _ = z.Update(msgSeason)

	return z
}

// Init created a new zess model
func (z *Model) Init() tea.Cmd {
	return nil
}

// Update updates the zess model
func (z *Model) Update(msg tea.Msg) (view.View, tea.Cmd) {
	switch msg := msg.(type) {
	// New scan(s)
	case scanMsg:
		z.lastScanID = msg.lastScanID
		// Add new scans
		for _, newScan := range msg.scans {
			found := false
			for i, modelScan := range z.scans {
				if newScan.time.equal(modelScan.time) {
					z.scans[i].amount++
					// Check for maxWeekScans
					if z.scans[i].amount > z.maxWeekScans {
						z.maxWeekScans = modelScan.amount
					}

					found = true
					break
				}
			}

			if !found {
				z.scans = append(z.scans, newScan)
				// Check for maxWeekScans
				if newScan.amount > z.maxWeekScans {
					z.maxWeekScans = newScan.amount
				}
				// Make sure the array doesn't get too big
				if len(z.scans) > config.GetDefaultInt("tui.zess.weeks", 10) {
					z.scans = z.scans[:1]
				}
			}

			// Update seasonScans
			z.seasonScans += newScan.amount
		}

	// New season!
	// Update variables accordinly
	case seasonMsg:
		z.currentSeason = msg.start
		z.seasonScans = 0
		z.maxWeekScans = 0

		validScans := make([]weekScan, 0, len(z.scans))

		for _, scan := range z.scans {
			// Add scans if they happend after (or in the same week of) the season start
			if scan.time.equal(z.currentSeason) || scan.time.after(z.currentSeason) {
				validScans = append(validScans, scan)

				if scan.amount > z.maxWeekScans {
					z.maxWeekScans = scan.amount
				}

				z.seasonScans += scan.amount
			}
		}

		z.scans = validScans
	}

	return z, nil
}

// View returns the view for the zess model
func (z *Model) View() string {
	chart := barchart.New(20, 20)

	for _, scan := range z.scans {
		bar := barchart.BarData{
			Label: scan.label,
			Values: []barchart.BarValue{{
				Name:  scan.label,
				Value: float64(scan.amount),
				Style: lipgloss.NewStyle().Foreground(lipgloss.Color("21")),
			}},
		}

		chart.Push(bar)
	}

	chart.Draw()

	style := lipgloss.NewStyle().Height(20).Align(lipgloss.Bottom).Render(lipgloss.JoinVertical(lipgloss.Left,
		fmt.Sprintf("Season scans\n%d", z.seasonScans),
		fmt.Sprintf("Max scans in a week\n%d", z.maxWeekScans),
	))

	return lipgloss.JoinHorizontal(lipgloss.Top,
		chart.View(),
		style,
	)
}

// GetUpdateDatas returns all the update functions for the zess model
func (z *Model) GetUpdateDatas() []view.UpdateData {
	return []view.UpdateData{
		{
			Name:     "zess scans",
			View:     z,
			Update:   updateScans,
			Interval: config.GetDefaultInt("tui.zess.interval_scan_s", 60),
		},
		{
			Name:     "zess season",
			View:     z,
			Update:   updateSeason,
			Interval: config.GetDefaultInt("tui.zess.interval_season_s", 3600),
		},
	}
}

// Check for any new scans
func updateScans(db *db.DB, view view.View) (tea.Msg, error) {
	z := view.(*Model)
	lastScanID := z.lastScanID

	// Get new scans
	scans, err := db.Queries.GetAllScansSinceID(context.Background(), lastScanID)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows shouldn't be considered an error
			err = nil
		}
		return nil, err
	}

	// No new scans
	if len(scans) == 0 {
		return nil, nil
	}

	zessScanMsg := scanMsg{lastScanID: lastScanID, scans: make([]weekScan, 0)}

	// Add new scans to scan msg
	for _, newScan := range scans {
		yearNumber, weekNumber := newScan.ScanTime.ISOWeek()
		newTime := time{year: yearNumber, week: weekNumber}

		found := false
		for i, scan := range zessScanMsg.scans {
			if scan.time.equal(newTime) {
				zessScanMsg.scans[i].amount++
				found = true
				break
			}
		}

		if !found {
			zessScanMsg.scans = append(zessScanMsg.scans, weekScan{time: newTime, amount: 1, label: newScan.ScanTime.Format("02/01")})
		}

		// Update scan ID
		// Not necessarly the first or last entry in the scans slice
		if newScan.ID > zessScanMsg.lastScanID {
			zessScanMsg.lastScanID = newScan.ID
		}
	}

	return zessScanMsg, nil
}

// Check if a new season started
func updateSeason(db *db.DB, view view.View) (tea.Msg, error) {
	z := view.(*Model)

	season, err := db.Queries.GetSeasonCurrent(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows shouldn't be considered an error
			err = nil
		}
		return nil, err
	}

	// Check if we have a new season
	yearNumber, weekNumber := season.Start.ISOWeek()
	seasonStart := time{year: yearNumber, week: weekNumber}
	if z.currentSeason.equal(seasonStart) {
		// Same season
		return nil, nil
	}

	return seasonMsg{start: seasonStart}, nil
}

func (z *time) equal(z2 time) bool {
	return z.week == z2.week && z.year == z2.year
}
func (z *time) after(z2 time) bool {
	if z.year > z2.year {
		return true
	} else if z.year < z2.year {
		return false
	}

	return z.week > z2.week
}

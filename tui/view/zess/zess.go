// Package zess provides the functions to draw an overview of the zess scans on a TUI
package zess

import (
	"context"
	"errors"
	"math/rand/v2"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jackc/pgx/v5"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/pkg/config"
	"github.com/zeusWPI/scc/tui/view"
	"go.uber.org/zap"
)

// yearWeek is used to represent a date by it's year and week
type yearWeek struct {
	year int
	week int
}

type weekScan struct {
	time   yearWeek
	amount int64
	start  string // The date when the week starts
	color  string
}

// Model represents the Model for the zess view
type Model struct {
	db            *db.DB
	lastScanID    int32
	scans         []weekScan // Scans per week
	showWeeks     int        // Amount of weeks to show
	maxWeekScans  int64
	currentSeason yearWeek // Start week of the season
	seasonScans   int64

	width  int
	height int
}

// Msg is the base message to indicate that something changed in the zess view
type Msg struct{}

// scanMsg is used to indicate that the zess view should be updated with new scans
type scanMsg struct {
	Msg
	lastScanID int32
	scans      []weekScan
}

// seasonMsg is used to indicate that the current season changed.
type seasonMsg struct {
	Msg
	start yearWeek
}

// NewModel creates a new zess model view
func NewModel(db *db.DB) view.View {
	m := &Model{
		db:            db,
		lastScanID:    -1,
		scans:         make([]weekScan, 0),
		showWeeks:     config.GetDefaultInt("tui.view.zess.weeks", 10),
		maxWeekScans:  -1,
		currentSeason: yearWeek{year: -1, week: -1},
		seasonScans:   0,
	}

	// Populate with data
	// The order in which this is called is important!
	msgScans, err := updateScans(m)
	if err != nil {
		zap.S().Error("TUI: Unable to update zess scans\n", err)
		return m
	}
	_, _ = m.Update(msgScans)

	msgSeason, err := updateSeason(m)
	if err != nil {
		zap.S().Error("TUI: Unable to update zess seasons\n", err)
		return m
	}
	_, _ = m.Update(msgSeason)

	return m
}

// Init created a new zess model
func (m *Model) Init() tea.Cmd {
	return nil
}

// Name returns the name of the view
func (m *Model) Name() string {
	return "Zess"
}

// Update updates the zess model
func (m *Model) Update(msg tea.Msg) (view.View, tea.Cmd) {
	switch msg := msg.(type) {
	case view.MsgSize:
		// Size update!
		// Check if it's relevant for this view
		entry, ok := msg.Sizes[m.Name()]
		if ok {
			// Update all dependent styles
			m.width = entry.Width
			m.height = entry.Height

			m.updateStyles()
		}

		return m, nil

	// New scan(s)
	case scanMsg:
		m.lastScanID = msg.lastScanID
		// Add new scans
		for _, newScan := range msg.scans {
			found := false
			for i, modelScan := range m.scans {
				if newScan.time.equal(modelScan.time) {
					m.scans[i].amount++
					// Check for maxWeekScans
					if m.scans[i].amount > m.maxWeekScans {
						m.maxWeekScans = modelScan.amount
					}

					found = true
					break
				}
			}

			if !found {
				m.scans = append(m.scans, newScan)
				// Check for maxWeekScans
				if newScan.amount > m.maxWeekScans {
					m.maxWeekScans = newScan.amount
				}
				// Make sure the array doesn't get too big
				if len(m.scans) > m.showWeeks {
					m.scans = m.scans[:1]
				}
			}

			// Update seasonScans
			m.seasonScans += newScan.amount
		}

	// New season!
	// Update variables accordinly
	case seasonMsg:
		m.currentSeason = msg.start
		m.seasonScans = 0
		m.maxWeekScans = 0

		validScans := make([]weekScan, 0, len(m.scans))

		for _, scan := range m.scans {
			// Add scans if they happend after (or in the same week of) the season start
			if scan.time.equal(m.currentSeason) || scan.time.after(m.currentSeason) {
				validScans = append(validScans, scan)

				if scan.amount > m.maxWeekScans {
					m.maxWeekScans = scan.amount
				}

				m.seasonScans += scan.amount
			}
		}

		m.scans = validScans
	}

	return m, nil
}

// View returns the view for the zess model
func (m *Model) View() string {
	chart := m.viewChart()
	stats := m.viewStats()

	// Join them together
	view := lipgloss.JoinHorizontal(lipgloss.Top, chart, stats)
	return view
}

// GetUpdateDatas returns all the update functions for the zess model
func (m *Model) GetUpdateDatas() []view.UpdateData {
	return []view.UpdateData{
		{
			Name:     "zess scans",
			View:     m,
			Update:   updateScans,
			Interval: config.GetDefaultInt("tui.view.zess.interval_scan_s", 60),
		},
		{
			Name:     "zess season",
			View:     m,
			Update:   updateSeason,
			Interval: config.GetDefaultInt("tui.view.zess.interval_season_s", 3600),
		},
	}
}

// Check for any new scans
func updateScans(view view.View) (tea.Msg, error) {
	m := view.(*Model)
	lastScanID := m.lastScanID

	// Get new scans
	scans, err := m.db.Queries.GetAllScansSinceID(context.Background(), lastScanID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
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
		yearNumber, weekNumber := newScan.ScanTime.Time.ISOWeek()
		newTime := yearWeek{year: yearNumber, week: weekNumber}

		found := false
		for i, scan := range zessScanMsg.scans {
			if scan.time.equal(newTime) {
				zessScanMsg.scans[i].amount++
				found = true
				break
			}
		}

		if !found {
			zessScanMsg.scans = append(zessScanMsg.scans, weekScan{
				time:   newTime,
				amount: 1,
				start:  newScan.ScanTime.Time.Format("02/01"),
				color:  randomColor(),
			})
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
func updateSeason(view view.View) (tea.Msg, error) {
	m := view.(*Model)

	season, err := m.db.Queries.GetSeasonCurrent(context.Background())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// No rows shouldn't be considered an error
			err = nil
		}
		return nil, err
	}

	// Check if we have a new season
	yearNumber, weekNumber := season.Start.Time.ISOWeek()
	seasonStart := yearWeek{year: yearNumber, week: weekNumber}
	if m.currentSeason.equal(seasonStart) {
		// Same season
		return nil, nil
	}

	return seasonMsg{start: seasonStart}, nil
}

func (z *yearWeek) equal(z2 yearWeek) bool {
	return z.week == z2.week && z.year == z2.year
}

func (z *yearWeek) after(z2 yearWeek) bool {
	if z.year > z2.year {
		return true
	} else if z.year < z2.year {
		return false
	}

	return z.week > z2.week
}

func randomColor() string {
	return colors[rand.IntN(len(colors))]
}

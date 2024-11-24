package view

import (
	"context"
	"database/sql"
	"time"

	"github.com/NimbleMarkets/ntcharts/canvas"
	"github.com/NimbleMarkets/ntcharts/linechart/wavelinechart"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"go.uber.org/zap"
)

// ZessModel represents the model for the zess view
type ZessModel struct {
	db         *db.DB
	lastScanID int64
	scans      []zessDayScan
	totalScans int64
}

type zessDayScan struct {
	date   time.Time
	amount int64
}

// ZessScanMsg represents the message to update the zess view
type ZessScanMsg struct {
	lastScanID int64
	scans      []zessDayScan
}

// ZessSeasonMsg represents the message to update the zess view
type ZessSeasonMsg struct {
	valid  bool
	amount int64
}

// NewZessModel creates a new zess model view
func NewZessModel(db *db.DB) *ZessModel {
	return &ZessModel{db: db, lastScanID: -1, scans: make([]zessDayScan, 0), totalScans: 0}
}

// Init created a new zess model
func (z *ZessModel) Init() tea.Cmd {
	return tea.Batch(updateScans(z.db, z.lastScanID), updateSeason(z.db))
}

// Update updates the zess model
func (z *ZessModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ZessScanMsg:
		z.lastScanID = msg.lastScanID
		for _, scan := range msg.scans {
			if len(z.scans) == 0 || scan.date.After(z.scans[len(z.scans)-1].date) {
				z.scans = append(z.scans, scan)
				// TODO: Potentially remove first element (scans = scans[1:])
				continue
			}

			for i := len(z.scans) - 1; i >= 0; i-- {
				if scan.date.Equal(z.scans[i].date) {
					z.scans[i].amount += scan.amount
					break
				}
			}
		}

		return z, updateScans(z.db, z.lastScanID)

	case ZessSeasonMsg:
		if msg.valid {
			z.totalScans = msg.amount
		}

		return z, updateSeason(z.db)
	}

	return z, nil
}

// View returns the view for the zess model
func (z *ZessModel) View() string {
	chart := wavelinechart.New(40, 20, wavelinechart.WithYRange(-2, 30))
	chart.XLabelFormatter = func(_ int, v float64) string {
		return time.Now().Add(-time.Duration(v*24) * time.Hour).Format("02")
	}

	now := time.Now().Truncate(24 * time.Hour)
	for _, scan := range z.scans {
		chart.Plot(canvas.Float64Point{X: now.Sub(scan.date).Hours() / 24, Y: float64(scan.amount)})
	}
	chart.Draw()

	return chart.View()
}

func updateScans(db *db.DB, lastScanID int64) tea.Cmd {
	return tea.Tick(60*time.Second, func(_ time.Time) tea.Msg {
		scan, err := db.Queries.GetLastScan(context.Background())
		if err != nil {
			if err != sql.ErrNoRows {
				zap.S().Error("DB: Failed to get last scan", err)
			}
			return ZessScanMsg{lastScanID: lastScanID, scans: []zessDayScan{}}
		}

		if scan.ID <= lastScanID {
			return ZessScanMsg{lastScanID: lastScanID, scans: []zessDayScan{}}
		}

		scans, err := db.Queries.GetAllScansSinceID(context.Background(), lastScanID)
		if err != nil {
			if err != sql.ErrNoRows {
				zap.S().Error("DB: Failed to get scan count by day", err)
			}
			return ZessScanMsg{lastScanID: lastScanID, scans: []zessDayScan{}}
		}

		zessMsg := ZessScanMsg{lastScanID: scan.ID, scans: []zessDayScan{}}
		for _, scan := range scans {
			date := scan.ScanTime.Truncate(24 * time.Hour)

			if len(zessMsg.scans) > 0 && zessMsg.scans[len(zessMsg.scans)-1].date.Equal(date) {
				// Already entry for that day
				zessMsg.scans[len(zessMsg.scans)-1].amount++
			} else {
				// New day entry
				zessMsg.scans = append(zessMsg.scans, zessDayScan{
					date:   date,
					amount: 1,
				})
			}
		}

		return zessMsg
	})
}

func updateSeason(db *db.DB) tea.Cmd {
	return tea.Tick(3600*time.Second, func(_ time.Time) tea.Msg {
		amount, err := db.Queries.GetScansInCurrentSeason(context.Background())
		if err != nil {
			if err != sql.ErrNoRows {
				zap.S().Error("DB: Failed to get scans in current season", err)
			}
			return ZessSeasonMsg{valid: false, amount: 0}
		}

		return ZessSeasonMsg{valid: true, amount: amount}
	})
}

package cmd

import (
	"time"

	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/internal/pkg/zess"
	"github.com/zeusWPI/scc/pkg/config"
	"go.uber.org/zap"
)

// Zess starts the zess instance
func Zess(db *db.DB) (*zess.Zess, chan bool, chan bool) {
	zess := zess.New(db)

	doneSeason := make(chan bool)
	intervalSeason := config.GetDefaultInt("backend.zess.interval_season_s", 300)

	doneScan := make(chan bool)
	intervalScan := config.GetDefaultInt("backend.zess.interval_scan_s", 60)

	go zessPeriodicSeasonUpdate(zess, doneSeason, intervalSeason)
	go zessPeriodicScanUpdate(zess, doneScan, intervalScan)

	return zess, doneSeason, doneScan
}

func zessPeriodicSeasonUpdate(zess *zess.Zess, done chan bool, interval int) {
	zap.S().Info("Zess: Starting periodic season update with an interval of ", interval, " seconds")

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// Run immediatly once
	zap.S().Info("Zess: Updating seasons")
	if err := zess.UpdateSeasons(); err != nil {
		zap.S().Error("Zess: Error updating seasons\n", err)
	}

	for {
		select {
		case <-done:
			zap.S().Info("Zess: Stopping periodic season update")
			return
		case <-ticker.C:
			// Update seasons
			zap.S().Info("Zess: Updating seasons")
			if err := zess.UpdateSeasons(); err != nil {
				zap.S().Error("Zess: Error updating seasons\n", err)
			}
		}
	}
}

func zessPeriodicScanUpdate(zess *zess.Zess, done chan bool, interval int) {
	zap.S().Info("Zess: Starting periodic scan update with an interval of ", interval, " seconds")

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// Run immediatly once
	zap.S().Info("Zess: Updating scans")
	if err := zess.UpdateScans(); err != nil {
		zap.S().Error("Zess: Error updating scans\n", err)
	}

	for {
		select {
		case <-done:
			zap.S().Info("Zess: Stopping periodic scan update")
			return
		case <-ticker.C:
			// Update scans
			zap.S().Info("Zess: Updating scans")
			if err := zess.UpdateScans(); err != nil {
				zap.S().Error("Zess: Error updating scans\n", err)
			}
		}
	}
}

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
	doneScan := make(chan bool)

	go zessPeriodicSeasonUpdate(zess, doneSeason)
	go zessPeriodicScanUpdate(zess, doneScan)

	return zess, doneSeason, doneScan
}

func zessPeriodicSeasonUpdate(zess *zess.Zess, done chan bool) {
	interval := config.GetDefaultInt("zess.interval_season_s", 300)
	zap.S().Info("Zess: Starting periodic season update with an interval of ", interval, " seconds")

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			zap.S().Info("Zess: Stopping periodic season update")
			return
		case <-ticker.C:
			// Update seasons
			zap.S().Info("Zess: Updating seasons")
			err := zess.UpdateSeasons()
			if err != nil {
				zap.S().Error("Zess: Error updating seasons\n", err)
			}
		}
	}
}

func zessPeriodicScanUpdate(zess *zess.Zess, done chan bool) {
	interval := config.GetDefaultInt("zess.interval_scan_s", 60)
	zap.S().Info("Zess: Starting periodic scan update with an interval of ", interval, " seconds")

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			zap.S().Info("Zess: Stopping periodic scan update")
			return
		case <-ticker.C:
			// Update scans
			zap.S().Info("Zess: Updating scans")
			err := zess.UpdateScans()
			if err != nil {
				zap.S().Error("Zess: Error updating scans\n", err)
			}
		}
	}
}

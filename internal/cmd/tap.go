package cmd

import (
	"time"

	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/internal/pkg/tap"
	"github.com/zeusWPI/scc/pkg/config"
	"go.uber.org/zap"
)

// Tap starts the tap instance
func Tap(db *db.DB) (*tap.Tap, chan bool) {
	tap := tap.New(db)
	done := make(chan bool)

	go tapPeriodicUpdate(tap, done)

	return tap, done
}

func tapPeriodicUpdate(tap *tap.Tap, done chan bool) {
	interval := config.GetDefaultInt("tap.interval_s", 60)
	zap.S().Info("Tap: Starting periodic update with an interval of ", interval, " seconds")

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			zap.S().Info("Tap: Stopping periodic update")
			return
		case <-ticker.C:
			// Update tap
			zap.S().Info("Tap: Updating tap")
			err := tap.Update()
			if err != nil {
				zap.S().Error("Tap: Error updating tap\n", err)
			}
		}
	}

}

package cmd

import (
	"context"
	"time"

	"github.com/zeusWPI/scc/internal/database/repository"
	"github.com/zeusWPI/scc/internal/tap"
	"github.com/zeusWPI/scc/pkg/config"
	"go.uber.org/zap"
)

// Tap starts the tap instance
func Tap(repo repository.Repository) (*tap.Tap, chan bool) {
	tap := tap.New(repo)
	done := make(chan bool)
	interval := config.GetDefaultInt("backend.tap.interval_s", 60)

	go tapPeriodicUpdate(tap, done, interval)

	return tap, done
}

// TODO: Figure out the context situation
func tapPeriodicUpdate(tap *tap.Tap, done chan bool, interval int) {
	zap.S().Info("Tap: Starting periodic update with an interval of ", interval, " seconds")

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// Run immediatly once
	zap.S().Info("Tap: Updating tap")
	if err := tap.Update(context.Background()); err != nil {
		zap.S().Error("Tap: Error updating tap\n", err)
	}

	for {
		select {
		case <-done:
			zap.S().Info("Tap: Stopping periodic update")
			return
		case <-ticker.C:
			// Update tap
			zap.S().Info("Tap: Updating tap")
			if err := tap.Update(context.Background()); err != nil {
				zap.S().Error("Tap: Error updating tap\n", err)
			}
		}
	}
}

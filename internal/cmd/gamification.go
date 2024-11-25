package cmd

import (
	"time"

	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/internal/pkg/gamification"
	"github.com/zeusWPI/scc/pkg/config"
	"go.uber.org/zap"
)

// Gamification starts the gamification instance
func Gamification(db *db.DB) (*gamification.Gamification, chan bool) {
	gam := gamification.New(db)
	done := make(chan bool)

	go gamificationPeriodicUpdate(gam, done)

	return gam, done
}

func gamificationPeriodicUpdate(gam *gamification.Gamification, done chan bool) {
	interval := config.GetDefaultInt("gamification.interval_s", 3600)
	zap.S().Info("Gamification: Starting periodic leaderboard update with an interval of ", interval, " seconds")

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// Run immediatly once
	zap.S().Info("Gamification: Updating leaderboard")
	err := gam.Update()
	if err != nil {
		zap.S().Error("gamification: Error updating leaderboard\n", err)
	}

	for {
		select {
		case <-done:
			zap.S().Info("Gamification: Stopping periodic leaderboard update")
			return
		case <-ticker.C:
			// Update leaderboard
			zap.S().Info("Gamification: Updating leaderboard")
			err := gam.Update()
			if err != nil {
				zap.S().Error("gamification: Error updating leaderboard\n", err)
			}
		}
	}
}

package cmd

import (
	"time"

	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/internal/pkg/event"
	"github.com/zeusWPI/scc/pkg/config"
	"go.uber.org/zap"
)

// Event starts the event instance
func Event(db *db.DB) (*event.Event, chan bool) {
	ev := event.New(db)
	done := make(chan bool)
	interval := config.GetDefaultInt("backend.event.interval_s", 3600)

	go eventPeriodicUpdate(ev, done, interval)

	return ev, done
}

func eventPeriodicUpdate(ev *event.Event, done chan bool, interval int) {
	zap.S().Info("Event: Starting periodic leaderboard update with an interval of ", interval, " seconds")

	ticker := time.NewTimer(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// Run immediatly once
	zap.S().Info("Event: Updating events")
	if err := ev.Update(); err != nil {
		zap.S().Error("Event: Error updating events\n", err)
	}

	for {
		select {
		case <-done:
			zap.S().Info("Event: Stopping periodic leaderboard update")
			return
		case <-ticker.C:
			// Update leaderboard
			zap.S().Info("Event: Updating events")
			if err := ev.Update(); err != nil {
				zap.S().Error("Event: Error updating events\n", err)
			}
		}
	}
}

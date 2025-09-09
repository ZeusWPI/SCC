package utils

import (
	"context"
	"time"

	"go.uber.org/zap"
)

func Periodic(name string, interval time.Duration, fn func(ctx context.Context) error, done chan bool) {
	zap.S().Infof("Starting periodic task for %s", name)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	update := func() {
		zap.S().Infof("Running %s", name)
		if err := fn(ctx); err != nil {
			zap.S().Errorf("Error %s | %v", name, err)
		}
	}

	// Run immediatly once
	update()

loop:
	for {
		select {
		case <-done:
			break loop
		case <-ticker.C:
			update()
		}
	}

	zap.S().Infof("Stopping periodic task for %s", name)
}

package cmd

import (
	"time"

	"github.com/zeusWPI/scc/internal/database/repository"
	"github.com/zeusWPI/scc/internal/tap"
	"github.com/zeusWPI/scc/pkg/config"
	"github.com/zeusWPI/scc/pkg/utils"
)

// Tap starts the tap instance
func Tap(repo repository.Repository) (*tap.Tap, chan bool) {
	tap := tap.New(repo)

	done := make(chan bool)
	interval := config.GetDefaultInt("backend.tap.interval_s", 60)

	go utils.Periodic(
		"Tap",
		time.Duration(interval)*time.Second,
		tap.Update,
		done,
	)

	return tap, done
}

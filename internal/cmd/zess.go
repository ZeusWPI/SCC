package cmd

import (
	"time"

	"github.com/zeusWPI/scc/internal/database/repository"
	"github.com/zeusWPI/scc/internal/zess"
	"github.com/zeusWPI/scc/pkg/config"
	"github.com/zeusWPI/scc/pkg/utils"
)

// Zess starts the zess instance
func Zess(repo repository.Repository) (*zess.Zess, chan bool) {
	zess := zess.New(repo)

	done := make(chan bool)
	intervalSeason := config.GetDefaultInt("backend.zess.interval_season_s", 300)
	intervalScan := config.GetDefaultInt("backend.zess.interval_scan_s", 60)

	go utils.Periodic(
		"Zess season",
		time.Duration(intervalSeason)*time.Second,
		zess.UpdateSeasons,
		done,
	)

	go utils.Periodic(
		"Zess scans",
		time.Duration(intervalScan)*time.Second,
		zess.UpdateScans,
		done,
	)

	return zess, done
}

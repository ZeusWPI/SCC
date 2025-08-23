// Main entry point for the tui
package main

import (
	"flag"
	"fmt"

	"github.com/zeusWPI/scc/internal/cmd"
	"github.com/zeusWPI/scc/internal/database/repository"
	"github.com/zeusWPI/scc/pkg/config"
	"github.com/zeusWPI/scc/pkg/db"
	"github.com/zeusWPI/scc/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// Config
	err := config.Init()
	if err != nil {
		panic(err)
	}

	screen := flag.String("screen", "", "TUI screen to start")
	flag.Parse()

	if *screen == "" {
		flag.PrintDefaults()
		return
	}

	// Logger
	zapLogger, err := logger.New(fmt.Sprintf("%s.log", *screen), false)
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(zapLogger)

	zap.S().Info("Initializing TUI")

	// Database
	db, err := db.NewPSQL()
	if err != nil {
		zap.S().Fatal("DB: Fatal error\n", err)
	}

	repo := repository.New(db)

	// TUI
	err = cmd.TUI(*repo, *screen)
	if err != nil {
		zap.S().Fatal("TUI: Fatal error\n", err)
	}
}

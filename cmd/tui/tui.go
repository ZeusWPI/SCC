// Main entry point for the tui
package main

import (
	"github.com/zeusWPI/scc/internal/cmd"
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

	// Logger
	zapLogger, err := logger.New("tui", false)
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

	// TUI
	err = cmd.TUI(db)
	if err != nil {
		zap.S().Fatal("TUI: Fatal error\n", err)
	}
}

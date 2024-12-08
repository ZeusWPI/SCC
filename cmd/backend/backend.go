// Main entry point for the backend
package main

import (
	"github.com/zeusWPI/scc/internal/cmd"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/pkg/config"
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
	zapLogger, err := logger.New("backend", true)
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(zapLogger)

	zap.S().Info("Initializing backend")

	// Database
	db, err := db.New()
	if err != nil {
		zap.S().Fatal("DB: Fatal error\n", err)
	}

	// // Tap
	// _, _ = cmd.Tap(db)

	// Zess
	_, _, _ = cmd.Zess(db)

	// // Gamification
	// _, _ = cmd.Gamification(db)

	// // Event
	// _, _ = cmd.Event(db)

	// Spotify
	spotify, err := cmd.Song(db)
	if err != nil {
		zap.S().Error("Spotify: Initiating error, integration will not work.\n", err)
	}

	// API
	cmd.API(db, spotify)
}

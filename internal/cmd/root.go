// Package cmd package contains all the commands to start the application
package cmd

import (
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/internal/pkg/spotify"
	"github.com/zeusWPI/scc/internal/pkg/tap"
	"github.com/zeusWPI/scc/pkg/config"
	"github.com/zeusWPI/scc/pkg/logger"
	"go.uber.org/zap"
)

// Execute starts the entire application
func Execute() {
	zap.S().Info("Initializing application")

	// Config
	err := config.Init()
	if err != nil {
		panic(err)
	}

	// Logger
	zapLogger, err := logger.New()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(zapLogger)

	// Database
	db, err := db.New()
	if err != nil {
		zap.S().Fatal("DB: Fatal error\n", err)
	}

	// Spotify
	spotify, err := spotify.New(db)
	if err != nil {
		zap.S().Error("Spotify: Initiating error, integration will not work.\n", err)
	}

	// Tap
	tap := tap.New(db)

	// Everything that needs to be initialized is done
	// Time to start all parts of the application
	zap.S().Info("Starting application")

	_ = tapCmd(tap)

	apiCmd(db, spotify)
}

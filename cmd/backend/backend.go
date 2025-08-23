// Main entry point for the backend
package main

import (
	"github.com/zeusWPI/scc/internal/cmd"
	"github.com/zeusWPI/scc/internal/database/repository"
	"github.com/zeusWPI/scc/internal/server"
	"github.com/zeusWPI/scc/internal/server/service"
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
	zapLogger, err := logger.New("backend", true)
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(zapLogger)

	zap.S().Info("Initializing backend")

	// Database
	db, err := db.NewPSQL()
	if err != nil {
		zap.S().Fatal("DB: Fatal error\n", err)
	}

	// Repository
	repo := repository.New(db)

	// Tap
	_, _ = cmd.Tap(db)

	// Zess
	_, _, _ = cmd.Zess(db)

	// Spotify
	spotify, err := cmd.Song(db)
	if err != nil {
		zap.S().Error("Spotify: Initiating error, integration will not work.\n", err)
	}

	// API
	service := service.New(*repo)
	api := server.New(*service, db.Pool())

	zap.S().Infof("Server is running on %s", api.Addr)
	if err := api.Listen(api.Addr); err != nil {
		zap.S().Fatalf("Failure while running the server %v", err)
	}
}

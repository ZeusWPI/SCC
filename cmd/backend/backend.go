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
		zap.S().Fatalf("DB: Fatal error %v", err)
	}

	// Repository
	repo := repository.New(db)

	var dones []chan bool

	// Tap
	_, done := cmd.Tap(*repo)
	dones = append(dones, done)

	// Zess
	_, done = cmd.Zess(*repo)
	dones = append(dones, done)

	// Song
	if err := cmd.Song(); err != nil {
		zap.S().Fatalf("Initialize song %v", err)
	}

	// API
	service := service.New(*repo)
	api := server.New(*service)

	zap.S().Infof("Server is running on %s", api.Addr)
	if err := api.Listen(api.Addr); err != nil {
		zap.S().Fatalf("Failure while running the server %v", err)
	}

	for _, done := range dones {
		done <- true
	}
}

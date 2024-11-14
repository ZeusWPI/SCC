// Package cmd package contains all the commands to start the application
package cmd

import (
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/pkg/config"
	"github.com/zeusWPI/scc/pkg/logger"
	"go.uber.org/zap"
)

// Execute starts the entire application
func Execute() {
	err := config.Init()
	if err != nil {
		panic(err)
	}

	zapLogger, err := logger.New()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(zapLogger)

	db, err := db.New()
	if err != nil {
		zap.S().Fatal("DB: Fatal error", err)
	}

	apiCmd(db)
}

package cmd

import (
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

	zapLogger := logger.New()
	zap.ReplaceGlobals(zapLogger)
}

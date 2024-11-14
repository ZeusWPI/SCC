// Package logger provides a logger instance
package logger

import (
	"os"

	"github.com/zeusWPI/scc/pkg/config"
	"go.uber.org/zap"
)

// New returns a new logger instance
func New() *zap.Logger {
	// Create logs directory
	os.Mkdir("logs", os.ModePerm)

	// Create logger
	var zapConfig zap.Config
	env := config.GetDefaultString("app.env", "development")
	if env == "development" {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}
	zapConfig.OutputPaths = []string{"logs/scc.log"}
	zapConfig.ErrorOutputPaths = []string{"logs/scc.log"}

	logger := zap.Must(zapConfig.Build())
	defer logger.Sync()

	return logger
}

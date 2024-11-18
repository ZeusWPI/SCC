// Package logger provides a logger instance
package logger

import (
	"fmt"
	"os"

	"github.com/zeusWPI/scc/pkg/config"
	"go.uber.org/zap"
)

// New returns a new logger instance
func New(logFile string) (*zap.Logger, error) {
	// Create logs directory
	err := os.Mkdir("logs", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	// Create logger
	var zapConfig zap.Config
	env := config.GetDefaultString("app.env", "development")
	if env == "development" {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}
	zapConfig.OutputPaths = []string{fmt.Sprintf("logs/%s.log", logFile)}
	zapConfig.ErrorOutputPaths = []string{fmt.Sprintf("logs/%s.log", logFile)}

	logger := zap.Must(zapConfig.Build())

	return logger, nil
}

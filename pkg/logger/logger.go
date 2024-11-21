// Package logger provides a logger instance
package logger

import (
	"fmt"
	"os"

	"github.com/zeusWPI/scc/pkg/config"
	"go.uber.org/zap"
)

// New returns a new logger instance
func New(logFile string, console bool) (*zap.Logger, error) {
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

	outputPaths := []string{fmt.Sprintf("logs/%s.log", logFile)}
	if console {
		outputPaths = append(outputPaths, "stdout")
	}

	errorOutputPaths := []string{fmt.Sprintf("logs/%s.log", logFile)}
	if console {
		errorOutputPaths = append(errorOutputPaths, "stderr")
	}

	zapConfig.OutputPaths = outputPaths
	zapConfig.ErrorOutputPaths = errorOutputPaths

	logger := zap.Must(zapConfig.Build())

	return logger, nil
}

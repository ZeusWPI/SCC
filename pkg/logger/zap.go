// Package logger initiates a zap logger
package logger

import (
	"fmt"
	"os"

	"github.com/zeusWPI/scc/pkg/config"
	"go.uber.org/zap"
)

func New(logFile string, console bool) (*zap.Logger, error) {
	err := os.Mkdir("logs", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return nil, fmt.Errorf("create logs directory %w", err)
	}

	outputPaths := []string{fmt.Sprintf("logs/%s.log", logFile)}
	if console {
		outputPaths = append(outputPaths, "stdout")
	}

	errorOutputPaths := []string{fmt.Sprintf("logs/%s.log", logFile)}
	if console {
		errorOutputPaths = append(errorOutputPaths, "stderr")
	}

	var logger *zap.Logger
	env := config.GetDefaultString("app.env", "development")

	if env != "production" {
		cfg := zap.NewDevelopmentConfig()
		cfg.OutputPaths = outputPaths
		cfg.ErrorOutputPaths = errorOutputPaths

		logger = zap.Must(cfg.Build(zap.AddStacktrace(zap.WarnLevel)))
	} else {
		cfg := zap.NewProductionConfig()
		cfg.Level.SetLevel(zap.WarnLevel)
		cfg.OutputPaths = outputPaths
		cfg.ErrorOutputPaths = errorOutputPaths

		logger = zap.Must(cfg.Build())
	}

	logger = logger.With(zap.String("env", env))

	return logger, nil
}

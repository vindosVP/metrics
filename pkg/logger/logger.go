// Package logger initializes and configures zap logger
package logger

import "go.uber.org/zap"

// Log consists configured logger instance
var Log *zap.Logger = zap.NewNop()

// Initialize initialize the logger
func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = zl
	return nil
}

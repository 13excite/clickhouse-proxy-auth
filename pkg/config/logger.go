// Package log contains methods and settings of the main logger
package config

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Defaults logger for the CLI mode
var (
	Base   = zap.NewNop()
	Logger = Base.Sugar()
)

// InitLogger loads a global logger for the service
func InitLogger(c *Config) error {
	logConfig := zap.NewProductionConfig()
	logConfig.Sampling = nil

	// Log Level
	var logLevel zapcore.Level
	if err := logLevel.Set(c.Level); err != nil {
		return fmt.Errorf("could not determine log level: %w", err)
	}
	logConfig.Level.SetLevel(logLevel)

	// Handle different logger encodings
	logConfig.Encoding = c.Encoding
	// Enable Color
	if c.Color {
		logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	logConfig.DisableStacktrace = c.DisableStacktrace
	// Use sane timestamp when logging to console
	if logConfig.Encoding == "console" {
		logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// JSON Fields
	logConfig.EncoderConfig.MessageKey = "msg"
	logConfig.EncoderConfig.LevelKey = "level"
	logConfig.EncoderConfig.CallerKey = "caller"

	// Settings
	logConfig.Development = c.DevMode
	logConfig.DisableCaller = c.DisableCaller

	logConfig.OutputPaths = c.OutputPaths
	logConfig.ErrorOutputPaths = c.ErrorOutputPaths

	// Build the logger
	globalLogger, err := logConfig.Build()
	if err != nil {
		return fmt.Errorf("could not build log config: %w", err)
	}
	zap.ReplaceGlobals(globalLogger)

	return nil
}

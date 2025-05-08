// Package logger provides a preconfigured zerolog-based logger with log level control via environment variables.
package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

const (
	// EnvLogLevel defines the environment variable name used to configure the logger level.
	EnvLogLevel = "LOG_LEVEL"
)

// InitLogger initializes and returns a zerolog.Logger configured for console output
// with log level determined by the LOG_LEVEL environment variable.
// Supported levels: debug, info, warn, error, fatal.
func InitLogger() zerolog.Logger {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    true,
		TimeFormat: "",
		PartsOrder: []string{
			zerolog.LevelFieldName,
			zerolog.MessageFieldName,
		},
	}
	return zerolog.New(output).Level(getLogLevel()).With().Logger()
}

// getLogLevel parses the LOG_LEVEL environment variable and returns the corresponding zerolog.Level.
// Defaults to info level if unset or unrecognized.
func getLogLevel() zerolog.Level {
	switch strings.ToLower(os.Getenv(EnvLogLevel)) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

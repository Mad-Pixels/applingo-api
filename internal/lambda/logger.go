package lambda

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

const (
	EnvLogLevel = "LOG_LEVEL"
)

func initLogger() zerolog.Logger {
	level := getLogLevel()

	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.MessageFieldName = "message"
	zerolog.LevelFieldName = "level"

	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339Nano,
		NoColor:    true,
	}
	logger := zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Caller().
		Logger()

	return logger
}

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

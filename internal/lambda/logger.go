package lambda

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

func initLogger() (*zap.Logger, error) {
	var (
		stacktrace bool
		develop    bool
		level      = zap.InfoLevel
	)
	switch strings.ToLower(os.Getenv(EnvLogLevel)) {
	case "debug":
		develop, stacktrace, level = true, true, zap.DebugLevel
	case "info":
		develop, stacktrace, level = true, true, zap.InfoLevel
	case "warn", "warning":
		develop, stacktrace, level = false, false, zap.WarnLevel
	case "error", "dpanic", "panic", "fatal":
		develop, stacktrace, level = false, false, zap.ErrorLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		EncodeLevel:    func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) { enc.AppendString(l.String()) },
		EncodeDuration: zapcore.StringDurationEncoder,
		MessageKey:     "message",
		LevelKey:       "level",
	}
	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		EncoderConfig:     encoderConfig,
		DisableStacktrace: !stacktrace,
		Development:       develop,
		Encoding:          "json",
		DisableCaller:     true,
		Sampling:          nil,
	}
	return config.Build()
}

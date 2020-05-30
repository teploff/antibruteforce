// Package logger provides logger functionality (using zap logger instead).
package logger

import (
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/teploff/antibruteforce/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	maxBackups = 3
	maxAge     = 3
)

type Option func(logger *zap.Logger) *zap.Logger

func New(dev bool, cfg *config.LoggerConfig, opts ...Option) *zap.Logger {
	var options []zap.Option

	prodConfig := zap.NewProductionEncoderConfig()
	prodConfig.TimeKey = "T"
	prodConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewConsoleEncoder(prodConfig)
	write := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize, // megabytes
		MaxBackups: maxBackups,  // old logs
		MaxAge:     maxAge,      // days
		Compress:   true,
	})

	if dev {
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		write = os.Stdout

		options = append(options, zap.AddStacktrace(zap.ErrorLevel))
		options = append(options, zap.Development())
	}

	core := zapcore.NewCore(
		encoder,
		write,
		getLogLevel(cfg.Level),
	)

	logger := zap.New(core, options...)
	for _, opt := range opts {
		logger = opt(logger)
	}

	return logger
}

// Unmarshal text to a zap level notation.
//
// level - text logging notation.
func getLogLevel(level string) zapcore.Level {
	lvl := zap.DebugLevel
	_ = lvl.UnmarshalText([]byte(level))

	return lvl
}

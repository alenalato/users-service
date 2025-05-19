package logger

import (
	"go.uber.org/zap/zapcore"
	"os"
	"strings"

	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func init() {
	logLevel := os.Getenv("LOG_LEVEL")
	var level zapcore.Level
	switch strings.ToLower(logLevel) {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	var config zap.Config
	if level == zap.DebugLevel {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}
	config.Level = zap.NewAtomicLevelAt(level)

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	Log = logger.Sugar()
}

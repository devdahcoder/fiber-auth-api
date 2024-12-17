package logger

import (
	"log/slog"
	"os"
)

var (
	globalSlogLogger *slog.Logger
)

type SlogLogConfig struct {
	Level slog.Level
	JSON  bool
}

func InitializeLogger(config ...SlogLogConfig) {

	var logHandler slog.Handler

	loggerConfig := SlogLogConfig{
		Level: slog.LevelInfo,
		JSON:  false,
	}

	if len(config) > 0 {
		loggerConfig = config[0]
	}

	if loggerConfig.JSON {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: loggerConfig.Level,
		})
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: loggerConfig.Level,
		})
	}

	globalLogger := slog.New(logHandler)

	slog.SetDefault(globalLogger)

}

func GetLogger() *slog.Logger {
	if globalSlogLogger == nil {
		InitializeLogger()
	}
	return globalSlogLogger
}

func Debug(msg string, args ...any) {
	globalSlogLogger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	globalSlogLogger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	globalSlogLogger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	globalSlogLogger.Error(msg, args...)
}

package logger

import (
	"log"
	"log/slog"
	"os"
)

type Logger struct {
	Log *slog.Logger // TODO
}

// Log
var logger *Logger

func GetLog() *Logger {
	return logger
}

func New(level string) *Logger {
	var l slog.Level
	switch level {
	case "debug":
		l = slog.LevelDebug
	case "info":
		l = slog.LevelInfo
	case "warn":
		l = slog.LevelWarn
	default:
		l = slog.LevelError
	}
	logger = &Logger{
		Log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: l})),
	}
	return logger
}

func (l Logger) Debug(msg string, args ...any) {
	l.Log.Debug(msg, args...)
}

func (l Logger) Info(msg string, args ...any) {
	l.Log.Info(msg, args...)
}

func (l Logger) Error(msg string, args ...any) {
	l.Log.Error(msg, args...)
}

func (l Logger) Warn(msg string, args ...any) {
	l.Log.Warn(msg, args...)
}

func (l Logger) Fatal(msg string) {
	log.Fatal(msg)
}

package logger

import (
	"log"
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger // TODO
}

// Log
var Log Logger

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
	return &Logger{
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: l})),
	}
}

func (l Logger) Debug(msg string, args ...any) {
	l.Debug(msg, args)
}

func (l Logger) Info(msg string, args ...any) {
	l.Info(msg, args)
}

func (l Logger) Error(msg string, args ...any) {
	l.Error(msg, args)
}

func (l Logger) Warn(msg string, args ...any) {
	l.Warn(msg, args)
}

func (l Logger) Fatal(msg string) {
	log.Fatal(msg)
}

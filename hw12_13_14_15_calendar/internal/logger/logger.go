package logger

import (
	"log"
	"log/slog"
	"os"
)

type Log struct {
	Log *slog.Logger
}

var logger *Log

func Logger() *Log {
	if logger == nil {
		logger = New("info")
	}
	return logger
}

func New(level string) *Log {
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
	logger = &Log{
		Log: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: l})),
	}
	return logger
}

func (l Log) Debug(msg string, args ...any) {
	l.Log.Debug(msg, args...)
}

func (l Log) Info(msg string, args ...any) {
	l.Log.Info(msg, args...)
}

func (l Log) Error(msg string, args ...any) {
	l.Log.Error(msg, args...)
}

func (l Log) Warn(msg string, args ...any) {
	l.Log.Warn(msg, args...)
}

func (l Log) Fatal(msg string) {
	log.Fatal(msg)
}

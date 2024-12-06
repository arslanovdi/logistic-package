// Package logger логирование
package logger

import (
	"log/slog"
	"os"
)

// InitializeLogger инициализирует slog логгера
func InitializeLogger(level slog.Level) {
	logger := slog.New(slog.NewJSONHandler(
		os.Stderr,
		&slog.HandlerOptions{
			Level: level,
		}))

	slog.SetDefault(logger)
}

package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Setup configures the default slog logger to write JSON logs to stdout with the specified log level.
func Setup(levelStr string) {
	level := ParseLevel(levelStr)

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

// ParseLevel converts a level string into a slog.Level.
func ParseLevel(levelStr string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(levelStr)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

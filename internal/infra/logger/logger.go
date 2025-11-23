package logger

import (
	"log/slog"
	"os"
)

func LoggerInit(env string) *slog.Logger {
	var handler slog.Handler

	if env == "development" || env == "dev" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}

package internal

import (
	"os"
	"log/slog"
)

func NewLogger(env string) *slog.Logger {
	level := slog.LevelInfo
	if env == "dev" {
		level = slog.LevelDebug
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: level,
    })

	return slog.New(handler)
}
package shared

import (
	"log/slog"
	"os"
)

func InitLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger) // biar bisa pakai slog.Info langsung

	return logger
}

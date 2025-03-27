package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

// ConfigureLogging will set a default logger to have colors
// Utilizes tint
func ConfigureLogging() {
	w := os.Stderr

	var level slog.Level
	if os.Getenv("ENV") == "dev" {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	// set global logger with custom options
	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      level,
			TimeFormat: time.Kitchen,
		}),
	))
}

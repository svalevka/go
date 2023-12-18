package common

import (
	"log/slog"
	"os"
)

// Logs contains common logging configuration used by all services.
type Logs struct {
	// Debug enabled debugging information in logs, otherwise discarded.
	Debug bool `yaml:"debug"`
}

// GetLogger returns a slog.Logger configured by Logs.
func (l *Logs) GetLogger() *slog.Logger {
	level := slog.LevelInfo
	if l.Debug {
		level = slog.LevelDebug
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
}

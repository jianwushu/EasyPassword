package logger

import (
	"io"
	"log/slog"
	"os"
)

// Init initializes the global logger with the given configuration.
func Init(level slog.Level, format string, output io.Writer) {
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: level,
	}

	if output == nil {
		output = os.Stdout
	}

	switch format {
	case "json":
		handler = slog.NewJSONHandler(output, opts)
	case "text":
		fallthrough
	default:
		handler = slog.NewTextHandler(output, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
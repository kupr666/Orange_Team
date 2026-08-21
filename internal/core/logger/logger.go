package core_logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

func New(config Config) (*slog.Logger, error) {
	return NewWithWriter(config, os.Stdout)
}

func NewWithWriter(
	config Config,
	output io.Writer,
) (*slog.Logger, error) {
	level, err := config.LogLevel()
	if err != nil {
		return nil, err
	}

	options := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}

	var handler slog.Handler

	switch config.Format {
	case "text":
		handler = slog.NewTextHandler(output, options)
	case "json":
		handler = slog.NewJSONHandler(output, options)
	default:
		return nil, fmt.Errorf(
			"unsupported logger format %q",
			config.Format,
		)
	}

	return slog.New(handler), nil
}
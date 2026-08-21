package core_logger

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Level  string `envconfig:"LEVEL" default:"INFO"`
	Format string `envconfig:"FORMAT" default:"text"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("LOGGER", &config); err != nil {
		return Config{}, fmt.Errorf("process logger config: %w", err)
	}

	return config, nil
}

func (c Config) LogLevel() (slog.Level, error) {
	switch strings.ToUpper(c.Level) {
	case "DEBUG":
		return slog.LevelDebug, nil
	case "INFO":
		return slog.LevelInfo, nil
	case "WARN":
		return slog.LevelWarn, nil
	case "ERROR":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("unsupported logger level %q", c.Level)
	}
}
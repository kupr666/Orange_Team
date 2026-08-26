package leaderboard_service

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Timezone         string        `envconfig:"TIMEZONE" default:"Europe/Moscow"`
	SnapshotInterval time.Duration `envconfig:"SNAPSHOT_INTERVAL" default:"1m"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("LEADERBOARD", &config); err != nil {
		return Config{}, fmt.Errorf("process leaderboard config: %w", err)
	}

	if config.SnapshotInterval <= 0 {
		return Config{}, fmt.Errorf("snapshot interval must be positive")
	}

	if _, err := time.LoadLocation(config.Timezone); err != nil {
		return Config{}, fmt.Errorf(
			"load leaderboard timezone %q: %w",
			config.Timezone,
			err,
		)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		panic(err)
	}
	return config
}

func (c Config) LocationMust() *time.Location {
	location, err := time.LoadLocation(c.Timezone)
	if err != nil {
		panic(fmt.Errorf("load leaderboard timezone %q: %w", c.Timezone, err))
	}
	return location
}

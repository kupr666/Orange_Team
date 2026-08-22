package core_logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Level  string `envconfig:"LEVEL"  default:"INFO"`
	Format string `envconfig:"FORMAT" default:"text"`
	Folder string `envconfig:"FOLDER" default:"logs"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("LOGGER", &config); err != nil {
		return Config{}, fmt.Errorf("process logger config: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get Logger config %w", err)
		panic(err)
	}

	return config
}

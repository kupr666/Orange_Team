package core_auth_jwt

import (
	"fmt"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

const minimumSecretLengthBytes = 32

type Config struct {
	Secret   string        `envconfig:"SECRET" required:"true"`
	Issuer   string        `envconfig:"ISSUER" required:"true"`
	Audience string        `envconfig:"AUDIENCE" required:"true"`
	TTL      time.Duration `envconfig:"TTL" required:"true"`
}

func NewConfig() (Config, error) { 
	var config Config 

	if err := envconfig.Process("JWT", &config); err != nil {
		return Config{}, fmt.Errorf("process JWT config: %w", err)
	}

	if err := config.Validate(); err != nil {
		return Config{}, fmt.Errorf("validate JWT config: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get JWT config: %w", err)
		panic(err)
	}

	return config
}

func (c Config) Validate() error {
	if len([]byte(c.Secret)) < minimumSecretLengthBytes {
		return fmt.Errorf(
			"JWT secret must contain at least %d bytes",
			minimumSecretLengthBytes,
		)
	}
	if strings.TrimSpace(c.Issuer) == "" {
		return fmt.Errorf("JWT issuer is empty")
	}
	if strings.TrimSpace(c.Audience) == "" {
		return fmt.Errorf("JWT audience is empty")
	}
	if c.TTL <= 0 {
		return fmt.Errorf("JWT TTL must be positive")
	}

	return nil
}

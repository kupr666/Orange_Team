package core_auth_jwt

import (
	"fmt"
	"strings"
	"time"
)

const minimumSecretLengthBytes = 32

type Config struct {
	Secret   string
	Issuer   string
	Audience string
	TTL      time.Duration
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

package core_auth_jwt

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	core_auth "github.com/kupr666/Orange_Team/internal/core/auth"
)

func TestManagerIssueAndVerifyAccessToken(t *testing.T) {
	manager := newTestManager(t, testConfig())
	now := time.Date(2026, time.August, 25, 12, 0, 0, 0, time.UTC)
	manager.now = func() time.Time { return now }

	want := core_auth.Principal{
		UserID: uuid.New(),
		Role:   "user",
	}

	token, err := manager.IssueAccessToken(want)
	if err != nil {
		t.Fatalf("IssueAccessToken() error = %v", err)
	}

	got, err := manager.VerifyAccessToken(token)
	if err != nil {
		t.Fatalf("VerifyAccessToken() error = %v", err)
	}
	if got != want {
		t.Fatalf("VerifyAccessToken() = %+v, want %+v", got, want)
	}
}

func TestManagerRejectsExpiredAccessToken(t *testing.T) {
	config := testConfig()
	manager := newTestManager(t, config)
	issuedAt := time.Date(2026, time.August, 25, 12, 0, 0, 0, time.UTC)
	manager.now = func() time.Time { return issuedAt }

	token, err := manager.IssueAccessToken(core_auth.Principal{
		UserID: uuid.New(),
		Role:   "user",
	})
	if err != nil {
		t.Fatalf("IssueAccessToken() error = %v", err)
	}

	manager.now = func() time.Time {
		return issuedAt.Add(config.TTL + time.Second)
	}
	if _, err := manager.VerifyAccessToken(token); err == nil {
		t.Fatal("VerifyAccessToken() accepted expired token")
	}
}

func TestManagerRejectsTokenSignedWithAnotherSecret(t *testing.T) {
	issuer := newTestManager(t, testConfig())
	verifierConfig := testConfig()
	verifierConfig.Secret = strings.Repeat("b", minimumSecretLengthBytes)
	verifier := newTestManager(t, verifierConfig)

	token, err := issuer.IssueAccessToken(core_auth.Principal{
		UserID: uuid.New(),
		Role:   "user",
	})
	if err != nil {
		t.Fatalf("IssueAccessToken() error = %v", err)
	}

	if _, err := verifier.VerifyAccessToken(token); err == nil {
		t.Fatal("VerifyAccessToken() accepted token signed with another secret")
	}
}

func TestManagerRejectsEmptyUserID(t *testing.T) {
	manager := newTestManager(t, testConfig())

	if _, err := manager.IssueAccessToken(core_auth.Principal{}); err == nil {
		t.Fatal("IssueAccessToken() accepted empty user ID")
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name   string
		change func(*Config)
	}{
		{
			name: "short secret",
			change: func(config *Config) {
				config.Secret = "short"
			},
		},
		{
			name: "empty issuer",
			change: func(config *Config) {
				config.Issuer = ""
			},
		},
		{
			name: "empty audience",
			change: func(config *Config) {
				config.Audience = ""
			},
		},
		{
			name: "non-positive TTL",
			change: func(config *Config) {
				config.TTL = 0
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := testConfig()
			test.change(&config)

			if _, err := NewManager(config); err == nil {
				t.Fatal("NewManager() accepted invalid config")
			}
		})
	}
}

func newTestManager(t *testing.T, config Config) *Manager {
	t.Helper()

	manager, err := NewManager(config)
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}

	return manager
}

func testConfig() Config {
	return Config{
		Secret:   strings.Repeat("a", minimumSecretLengthBytes),
		Issuer:   "solo-leveling-api",
		Audience: "solo-leveling-client",
		TTL:      15 * time.Minute,
	}
}

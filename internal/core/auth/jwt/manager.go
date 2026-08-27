package core_auth_jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	core_auth "github.com/kupr666/Orange_Team/internal/core/auth"
)

// обьект, который реализует интерфейс verifier (т.к у него есть)
// метод VerifyAccessToken(token string) (core_auth.Principal, error)
type Manager struct {
	secret   []byte
	issuer   string
	audience string
	ttl      time.Duration
	now      func() time.Time
}

func NewManager(config Config) (*Manager, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf(
			"validate JWT config: %w",
			err,
		)
	}

	return &Manager{
		secret:   []byte(config.Secret),
		issuer:   config.Issuer,
		audience: config.Audience,
		ttl:      config.TTL,
		now:      time.Now,
	}, nil
}

// creates new signed JWT for verified user
func (m *Manager) IssueAccessToken(
	principal core_auth.Principal,
) (string, error) {
	if principal.UserID == uuid.Nil {
		return "", fmt.Errorf("issue access token: user ID is empty")
	}

	now := m.now().UTC()
	tokenClaims := claims{
		Role: principal.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   principal.UserID.String(),
			Audience:  jwt.ClaimStrings{m.audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
	}

	signedToken, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		tokenClaims,
	).SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}

	return signedToken, nil
}

func (m *Manager) VerifyAccessToken(
	tokenString string,
) (core_auth.Principal, error) {
	if tokenString == "" {
		return core_auth.Principal{}, fmt.Errorf("access token is empty")
	}

	// prepare struct
	tokenClaims := &claims{}
	// move cliams to higher struct
	token, err := jwt.ParseWithClaims(
		tokenString,
		tokenClaims,
		func(_ *jwt.Token) (any, error) {
			return m.secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(m.issuer),
		jwt.WithAudience(m.audience),
		jwt.WithTimeFunc(m.now),
	)
	if err != nil {
		return core_auth.Principal{}, fmt.Errorf("parse access token: %w", err)
	}
	if !token.Valid {
		return core_auth.Principal{}, errors.New("access token is invalid")
	}

	// extract userID from claims
	userID, err := uuid.Parse(tokenClaims.Subject)
	if err != nil {
		return core_auth.Principal{}, fmt.Errorf(
			"parse access token subject: %w",
			err,
		)
	}
	if userID == uuid.Nil {
		return core_auth.Principal{}, errors.New(
			"access token subject contains empty user ID",
		)
	}

	return core_auth.Principal{
		UserID: userID,
		Role:   tokenClaims.Role,
	}, nil
}

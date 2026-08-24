package core_http_middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

type contextKey string

const userIDContextKey contextKey = "userID"

func JWT(secret string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := core_logger.FromContext(r.Context())

			tokenString, ok := bearerToken(
				r.Header.Get("Authorization"),
			)
			if !ok {
				log.Debug(
					"request rejected by authentication middleware",
					slog.String("reason", "bearer token is missing"),
				)

				writeUnauthorized(w)
				return
			}

			userID, err := parseJWT(tokenString, secret)
			if err != nil {
				log.Debug(
					"request rejected by authentication middleware",
					slog.String("reason", "invalid or expired token"),
					slog.Any("error", err),
				)

				writeUnauthorized(w)
				return
			}

			ctx := context.WithValue(
				r.Context(),
				userIDContextKey,
				userID,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func bearerToken(authorization string) (string, bool) {
	parts := strings.Fields(authorization)

	if len(parts) != 2 ||
		!strings.EqualFold(parts[0], "Bearer") ||
		parts[1] == "" {
		return "", false
	}

	return parts[1], true
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDContextKey).(uuid.UUID)
	return userID, ok
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("WWW-Authenticate", "Bearer")
	w.WriteHeader(http.StatusUnauthorized)

	_ = json.NewEncoder(w).Encode(core_http_response.ErrorResponse{
		Error:   "unauthorized",
		Message: "valid JWT token is required",
	})
}
func parseJWT(tokenString, secret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(secret), nil
		},
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse JWT: %w", err)
	}

	if !token.Valid {
		return uuid.Nil, errors.New("invalid JWT")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse JWT subject: %w", err)
	}

	return userID, nil
}

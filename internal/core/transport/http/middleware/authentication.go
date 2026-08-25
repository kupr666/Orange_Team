package core_http_middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	core_auth "github.com/kupr666/Orange_Team/internal/core/auth"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

const unauthorizedMessage = "valid JWT token is required"

type AccessTokenVerifier interface {
	VerifyAccessToken(token string) (core_auth.Principal, error)
}

func Authentication(verifier AccessTokenVerifier) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := core_logger.FromContext(r.Context())
			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

			token, ok := bearerToken(r.Header.Get("Authorization"))
			if !ok {
				respondAuthenticationFailure(
					responseHandler,
					errors.New("bearer token is missing or malformed"),
				)
				return
			}

			principal, err := verifier.VerifyAccessToken(token)
			if err != nil {
				respondAuthenticationFailure(responseHandler, err)
				return
			}

			ctx := core_auth.WithPrincipal(r.Context(), principal)
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

func respondAuthenticationFailure(
	responseHandler *core_http_response.HTTPResponseHandler,
	reason error,
) {
	responseHandler.ErrorResponse(
		fmt.Errorf(
			"authenticate request: %v: %w",
			reason,
			core_errors.ErrUnauthorized,
		),
		unauthorizedMessage,
	)
}

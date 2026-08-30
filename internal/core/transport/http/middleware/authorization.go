package core_http_middleware

import (
	"errors"
	"fmt"
	"net/http"

	core_auth "github.com/kupr666/Orange_Team/internal/core/auth"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

const forbiddenMessage = "insufficient permissions"

func RequireRole(allowedRoles ...string) Middleware {
	allowed := make(map[string]struct{}, len(allowedRoles))
	for _, role := range allowedRoles {
		allowed[role] = struct{}{}
	}
	

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := core_logger.FromContext(r.Context())
			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

			principal, ok := core_auth.PrincipalFromContext(r.Context())
			if !ok {
				respondAuthenticationFailure(
					responseHandler,
					errors.New("authenticated principal is missing"),
						
				)
				return
			}

			if _, ok := allowed[principal.Role]; !ok {
				responseHandler.ErrorResponse(
					fmt.Errorf(
						"role %q is not allowed: %w",
						principal.Role,
						core_errors.ErrForbidden,
					),
					forbiddenMessage,
				)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

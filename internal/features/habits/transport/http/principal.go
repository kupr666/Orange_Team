package habits_transport_http

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	core_auth "github.com/kupr666/Orange_Team/internal/core/auth"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
)

func authenticatedUserID(r *http.Request) (uuid.UUID, error) {
	principal, ok := core_auth.PrincipalFromContext(r.Context())
	if !ok {
		return uuid.Nil, fmt.Errorf("authenticated principal is missing: %w", core_errors.ErrUnauthorized)
	}
	return principal.UserID, nil
}

package core_http_middleware

import (
	"context"

	"github.com/google/uuid"
	core_auth "github.com/kupr666/Orange_Team/internal/core/auth"
)

// UserIDFromContext keeps workout handlers compatible with the shared
// authenticated Principal stored by Authentication middleware.
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	principal, ok := core_auth.PrincipalFromContext(ctx)
	if !ok || principal.UserID == uuid.Nil {
		return uuid.Nil, false
	}

	return principal.UserID, true
}

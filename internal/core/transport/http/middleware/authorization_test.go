package core_http_middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	core_auth "github.com/kupr666/Orange_Team/internal/core/auth"
)

func TestRequireRoleRejectsMissingPrincipal(t *testing.T) {
	handler := RequireRole("admin")(http.HandlerFunc(
		func(http.ResponseWriter, *http.Request) {
			t.Fatal("next handler was called without principal")
		},
	))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, requestWithLogger(t))

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
	}
}

func TestRequireRoleRejectsDisallowedRole(t *testing.T) {
	handler := RequireRole("admin")(http.HandlerFunc(
		func(http.ResponseWriter, *http.Request) {
			t.Fatal("next handler was called for forbidden role")
		},
	))
	request := requestWithLogger(t)
	ctx := core_auth.WithPrincipal(request.Context(), core_auth.Principal{
		UserID: uuid.New(),
		Role:   "user",
	})
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request.WithContext(ctx))

	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusForbidden)
	}
	if got := response.Header().Get("WWW-Authenticate"); got != "" {
		t.Fatalf("WWW-Authenticate = %q, want empty for 403", got)
	}
}

func TestRequireRoleAllowsMatchingRole(t *testing.T) {
	handler := RequireRole("admin", "moderator")(http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	))
	request := requestWithLogger(t)
	ctx := core_auth.WithPrincipal(request.Context(), core_auth.Principal{
		UserID: uuid.New(),
		Role:   "admin",
	})
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request.WithContext(ctx))

	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
	}
}

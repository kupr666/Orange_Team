package core_http_middleware

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	core_auth "github.com/kupr666/Orange_Team/internal/core/auth"
	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
)

type accessTokenVerifierStub struct {
	principal core_auth.Principal
	err       error
	token     string
	called    bool
}

func (v *accessTokenVerifierStub) VerifyAccessToken(
	token string,
) (core_auth.Principal, error) {
	v.called = true
	v.token = token

	return v.principal, v.err
}

func TestAuthenticationRejectsMissingBearerToken(t *testing.T) {
	verifier := &accessTokenVerifierStub{}
	nextCalled := false
	handler := Authentication(verifier)(http.HandlerFunc(
		func(http.ResponseWriter, *http.Request) {
			nextCalled = true
		},
	))

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, requestWithLogger(t))

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
	}
	if got := response.Header().Get("WWW-Authenticate"); got != "Bearer" {
		t.Fatalf("WWW-Authenticate = %q, want Bearer", got)
	}
	if verifier.called {
		t.Fatal("verifier was called without a bearer token")
	}
	if nextCalled {
		t.Fatal("next handler was called for unauthenticated request")
	}
}

func TestAuthenticationRejectsInvalidAccessToken(t *testing.T) {
	verifier := &accessTokenVerifierStub{err: errors.New("invalid signature")}
	handler := Authentication(verifier)(http.HandlerFunc(
		func(http.ResponseWriter, *http.Request) {
			t.Fatal("next handler was called for invalid token")
		},
	))
	request := requestWithLogger(t)
	request.Header.Set("Authorization", "Bearer invalid-token")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
	}
}

func TestAuthenticationAddsPrincipalToContext(t *testing.T) {
	want := core_auth.Principal{
		UserID: uuid.New(),
		Role:   "user",
	}
	verifier := &accessTokenVerifierStub{principal: want}
	handler := Authentication(verifier)(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			got, ok := core_auth.PrincipalFromContext(r.Context())
			if !ok {
				t.Fatal("principal is missing from request context")
			}
			if got != want {
				t.Fatalf("principal = %+v, want %+v", got, want)
			}
			w.WriteHeader(http.StatusNoContent)
		},
	))
	request := requestWithLogger(t)
	request.Header.Set("Authorization", "Bearer valid-token")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
	}
	if verifier.token != "valid-token" {
		t.Fatalf("verified token = %q, want valid-token", verifier.token)
	}
}

func requestWithLogger(t *testing.T) *http.Request {
	t.Helper()

	log := &core_logger.Logger{
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := core_logger.ToContext(request.Context(), log)

	return request.WithContext(ctx)
}

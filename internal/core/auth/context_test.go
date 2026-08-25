package core_auth

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestPrincipalContextRoundTrip(t *testing.T) {
	want := Principal{
		UserID: uuid.New(),
		Role:   "user",
	}

	ctx := WithPrincipal(context.Background(), want)
	got, ok := PrincipalFromContext(ctx)
	if !ok {
		t.Fatal("PrincipalFromContext() did not find principal")
	}
	if got != want {
		t.Fatalf("PrincipalFromContext() = %+v, want %+v", got, want)
	}
}

func TestPrincipalFromContextReturnsFalseWhenMissing(t *testing.T) {
	_, ok := PrincipalFromContext(context.Background())
	if ok {
		t.Fatal("PrincipalFromContext() found unexpected principal")
	}
}

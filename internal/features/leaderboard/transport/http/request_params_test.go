package leaderboard_transport_http

import (
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	leaderboard_service "github.com/kupr666/Orange_Team/internal/features/leaderboard/service"
)

func TestLeaderboardRequestParams(t *testing.T) {
	userID := uuid.New()
	request := httptest.NewRequest(
		"GET",
		"/leaderboard/daily?user_id="+userID.String()+"&limit=25",
		nil,
	)

	gotUserID, gotLimit, err := leaderboardRequestParams(request)
	if err != nil {
		t.Fatalf("parse leaderboard params: %v", err)
	}
	if gotUserID != userID {
		t.Fatalf("user ID = %s, want %s", gotUserID, userID)
	}
	if gotLimit != 25 {
		t.Fatalf("limit = %d, want 25", gotLimit)
	}
}

func TestLeaderboardRequestParamsUsesDefaultLimit(t *testing.T) {
	userID := uuid.New()
	request := httptest.NewRequest(
		"GET",
		"/leaderboard/daily?user_id="+userID.String(),
		nil,
	)

	_, gotLimit, err := leaderboardRequestParams(request)
	if err != nil {
		t.Fatalf("parse leaderboard params: %v", err)
	}
	if gotLimit != leaderboard_service.DefaultLimit {
		t.Fatalf(
			"limit = %d, want %d",
			gotLimit,
			leaderboard_service.DefaultLimit,
		)
	}
}

func TestLeaderboardRequestParamsRequiresUserID(t *testing.T) {
	request := httptest.NewRequest("GET", "/leaderboard/daily", nil)

	if _, _, err := leaderboardRequestParams(request); err == nil {
		t.Fatal("expected missing user_id error")
	}
}

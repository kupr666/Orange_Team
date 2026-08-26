package leaderboard_service

import (
	"testing"
	"time"
)

func TestLeaderboardWindows(t *testing.T) {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		t.Fatalf("load location: %v", err)
	}

	now := time.Date(2026, time.August, 26, 14, 30, 0, 0, location)

	tests := []struct {
		name      string
		window    periodWindow
		wantStart time.Time
		wantEnd   time.Time
	}{
		{
			name:   "daily",
			window: dailyWindow(now, location),
			wantStart: time.Date(
				2026, time.August, 26, 0, 0, 0, 0, location,
			),
			wantEnd: time.Date(
				2026, time.August, 27, 0, 0, 0, 0, location,
			),
		},
		{
			name:   "previous weekly",
			window: previousWeeklyWindow(now, location),
			wantStart: time.Date(
				2026, time.August, 17, 0, 0, 0, 0, location,
			),
			wantEnd: time.Date(
				2026, time.August, 24, 0, 0, 0, 0, location,
			),
		},
		{
			name:   "previous monthly",
			window: previousMonthlyWindow(now, location),
			wantStart: time.Date(
				2026, time.July, 1, 0, 0, 0, 0, location,
			),
			wantEnd: time.Date(
				2026, time.August, 1, 0, 0, 0, 0, location,
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !test.window.start.Equal(test.wantStart) {
				t.Fatalf(
					"start = %s, want %s",
					test.window.start,
					test.wantStart,
				)
			}
			if !test.window.end.Equal(test.wantEnd) {
				t.Fatalf(
					"end = %s, want %s",
					test.window.end,
					test.wantEnd,
				)
			}
		})
	}
}

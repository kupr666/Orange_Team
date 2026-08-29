package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHabitCompleteMaintainsDailyStreak(t *testing.T) {
	today := time.Date(2026, time.August, 30, 12, 0, 0, 0, time.UTC)
	habit := Habit{UserID: uuid.New(), Name: "Read", Description: "Ten pages"}

	first, changed := habit.Complete(today.AddDate(0, 0, -1))
	if !changed || first.CurrentStreak != 1 {
		t.Fatalf("first completion = (%d, %t), want (1, true)", first.CurrentStreak, changed)
	}

	second, changed := first.Complete(today)
	if !changed || second.CurrentStreak != 2 || !second.CompletedToday {
		t.Fatalf("second completion = (%d, %t, %t), want (2, true, true)", second.CurrentStreak, changed, second.CompletedToday)
	}

	duplicate, changed := second.Complete(today)
	if changed || duplicate.CurrentStreak != 2 {
		t.Fatalf("duplicate completion = (%d, %t), want (2, false)", duplicate.CurrentStreak, changed)
	}
}

func TestHabitCompleteResetsStreakAfterMissedDay(t *testing.T) {
	today := time.Date(2026, time.August, 30, 12, 0, 0, 0, time.UTC)
	lastCompleted := today.AddDate(0, 0, -2)
	habit := Habit{CurrentStreak: 8, LastCompletedDate: &lastCompleted}

	view := habit.ViewAt(today)
	if view.CurrentStreak != 0 {
		t.Fatalf("stale streak = %d, want 0", view.CurrentStreak)
	}

	completed, changed := habit.Complete(today)
	if !changed || completed.CurrentStreak != 1 {
		t.Fatalf("new streak = (%d, %t), want (1, true)", completed.CurrentStreak, changed)
	}
}

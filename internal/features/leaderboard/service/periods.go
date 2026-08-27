package leaderboard_service

import "time"

type periodWindow struct {
	start       time.Time
	end         time.Time
	nextRefresh time.Time
}

func dailyWindow(now time.Time, location *time.Location) periodWindow {
	localNow := now.In(location)
	start := time.Date(
		localNow.Year(),
		localNow.Month(),
		localNow.Day(),
		0, 0, 0, 0,
		location,
	)
	end := start.AddDate(0, 0, 1)

	return periodWindow{start: start, end: end, nextRefresh: end}
}

func previousWeeklyWindow(now time.Time, location *time.Location) periodWindow {
	localNow := now.In(location)
	today := time.Date(
		localNow.Year(),
		localNow.Month(),
		localNow.Day(),
		0, 0, 0, 0,
		location,
	)

	daysSinceMonday := (int(today.Weekday()) + 6) % 7
	currentWeekStart := today.AddDate(0, 0, -daysSinceMonday)
	previousWeekStart := currentWeekStart.AddDate(0, 0, -7)

	return periodWindow{
		start:       previousWeekStart,
		end:         currentWeekStart,
		nextRefresh: currentWeekStart.AddDate(0, 0, 7),
	}
}

func previousMonthlyWindow(now time.Time, location *time.Location) periodWindow {
	localNow := now.In(location)
	currentMonthStart := time.Date(
		localNow.Year(),
		localNow.Month(),
		1,
		0, 0, 0, 0,
		location,
	)
	previousMonthStart := currentMonthStart.AddDate(0, -1, 0)

	return periodWindow{
		start:       previousMonthStart,
		end:         currentMonthStart,
		nextRefresh: currentMonthStart.AddDate(0, 1, 0),
	}
}

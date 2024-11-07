package models

import "time"

func between(start, end time.Time, timestamp int64) bool {
	t := time.Unix(timestamp, 0)

	return (t.Equal(end) || t.After(end)) && (t.Equal(start) || t.Before(start))
}

func now() time.Time {
	now := time.Now()

	return time.Date(
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location(),
	)
}

func day(timestamp int64) time.Time {
	t := time.Unix(timestamp, 0)

	return time.Date(
		t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location(),
	)
}

func today() time.Time {
	now := time.Now()

	return time.Date(
		now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location(),
	)
}

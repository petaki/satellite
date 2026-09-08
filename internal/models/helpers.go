package models

import (
	"strings"
	"time"
)

func escapeGlob(value string) string {
	var builder strings.Builder

	builder.Grow(len(value))

	for _, r := range value {
		switch r {
		case '*', '?', '[', ']', '\\':
			builder.WriteRune('\\')
		}

		builder.WriteRune(r)
	}

	return builder.String()
}

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

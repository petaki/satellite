package models

import "time"

func today() time.Time {
	now := time.Now()

	return time.Date(
		now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location(),
	)
}

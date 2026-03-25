package dto

import "time"

func nullableTime(value time.Time) *time.Time {
	if value.IsZero() {
		return nil
	}

	utc := value.UTC()
	return &utc
}

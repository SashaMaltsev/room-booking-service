package postgres

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	scheduledomain "github.com/SashaMaltsev/room-booking-service/internal/domain/schedule"
)

func nullableString(value *string) any {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return trimmed
}

func nullableInt(value *int) any {
	if value == nil {
		return nil
	}

	return *value
}

func nullableTime(value *time.Time) any {
	if value == nil {
		return nil
	}

	return value.UTC()
}

func encodeWeekdays(days []scheduledomain.Weekday) string {
	if len(days) == 0 {
		return "{}"
	}

	parts := make([]string, 0, len(days))
	for _, day := range days {
		parts = append(parts, strconv.Itoa(int(day)))
	}

	return "{" + strings.Join(parts, ",") + "}"
}

func decodeWeekdays(raw string) ([]scheduledomain.Weekday, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || trimmed == "{}" {
		return nil, nil
	}

	trimmed = strings.TrimPrefix(trimmed, "{")
	trimmed = strings.TrimSuffix(trimmed, "}")
	if trimmed == "" {
		return nil, nil
	}

	parts := strings.Split(trimmed, ",")
	days := make([]scheduledomain.Weekday, 0, len(parts))

	for _, part := range parts {
		value, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return nil, fmt.Errorf("parse weekday %q: %w", part, err)
		}

		day, err := scheduledomain.ParseWeekday(value)
		if err != nil {
			return nil, err
		}

		days = append(days, day)
	}

	return days, nil
}

func utcDayBounds(date time.Time) (time.Time, time.Time) {
	day := date.UTC()
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	return start, start.Add(24 * time.Hour)
}

package postgres

import (
	"testing"
	"time"

	scheduledomain "github.com/SashaMaltsev/room-booking-service/internal/domain/schedule"
)

func TestEncodeDecodeWeekdays(t *testing.T) {
	original := []scheduledomain.Weekday{
		scheduledomain.Monday,
		scheduledomain.Wednesday,
		scheduledomain.Friday,
	}

	encoded := encodeWeekdays(original)
	if encoded != "{1,3,5}" {
		t.Fatalf("expected {1,3,5}, got %q", encoded)
	}

	decoded, err := decodeWeekdays(encoded)
	if err != nil {
		t.Fatalf("expected decode without error, got %v", err)
	}

	if len(decoded) != len(original) {
		t.Fatalf("expected %d weekdays, got %d", len(original), len(decoded))
	}

	for i, day := range original {
		if decoded[i] != day {
			t.Fatalf("expected weekday %v at index %d, got %v", day, i, decoded[i])
		}
	}
}

func TestUTCDateBounds(t *testing.T) {
	date := time.Date(2026, time.March, 24, 22, 30, 0, 0, time.FixedZone("UTC+3", 3*60*60))

	start, end := utcDayBounds(date)

	expectedStart := time.Date(2026, time.March, 24, 0, 0, 0, 0, time.UTC)
	expectedEnd := expectedStart.Add(24 * time.Hour)

	if !start.Equal(expectedStart) {
		t.Fatalf("expected start %s, got %s", expectedStart, start)
	}

	if !end.Equal(expectedEnd) {
		t.Fatalf("expected end %s, got %s", expectedEnd, end)
	}
}

package schedule

import (
	"errors"
	"testing"
	"time"
)

func TestScheduleNormalizeSortsAndDeduplicatesWeekdays(t *testing.T) {
	start, _ := NewTimeOfDay(9, 0)
	end, _ := NewTimeOfDay(11, 0)

	entity := Schedule{
		RoomID:     "room-1",
		DaysOfWeek: []Weekday{Friday, Monday, Friday, Wednesday},
		StartTime:  start,
		EndTime:    end,
	}

	entity.Normalize()

	expected := []Weekday{Monday, Wednesday, Friday}
	if len(entity.DaysOfWeek) != len(expected) {
		t.Fatalf("expected %d weekdays, got %d", len(expected), len(entity.DaysOfWeek))
	}

	for i, day := range expected {
		if entity.DaysOfWeek[i] != day {
			t.Fatalf("expected weekday %v at index %d, got %v", day, i, entity.DaysOfWeek[i])
		}
	}
}

func TestScheduleValidateRejectsMisalignedWindow(t *testing.T) {
	start, _ := NewTimeOfDay(9, 15)
	end, _ := NewTimeOfDay(10, 0)

	entity := Schedule{
		RoomID:     "room-1",
		DaysOfWeek: []Weekday{Monday},
		StartTime:  start,
		EndTime:    end,
	}

	err := entity.Validate()
	if !errors.Is(err, ErrWindowNotAlignedToSlotDuration) {
		t.Fatalf("expected ErrWindowNotAlignedToSlotDuration, got %v", err)
	}
}

func TestScheduleWindowsForDate(t *testing.T) {
	start, _ := NewTimeOfDay(9, 0)
	end, _ := NewTimeOfDay(11, 0)

	entity := Schedule{
		RoomID:     "room-1",
		DaysOfWeek: []Weekday{Monday},
		StartTime:  start,
		EndTime:    end,
	}

	date := time.Date(2026, time.March, 23, 15, 45, 0, 0, time.FixedZone("UTC+3", 3*60*60))
	windows := entity.WindowsForDate(date)

	if len(windows) != 4 {
		t.Fatalf("expected 4 windows, got %d", len(windows))
	}

	expectedStart := time.Date(2026, time.March, 23, 9, 0, 0, 0, time.UTC)
	expectedEnd := time.Date(2026, time.March, 23, 9, 30, 0, 0, time.UTC)
	if !windows[0].Start.Equal(expectedStart) {
		t.Fatalf("expected first window start %s, got %s", expectedStart, windows[0].Start)
	}

	if !windows[0].End.Equal(expectedEnd) {
		t.Fatalf("expected first window end %s, got %s", expectedEnd, windows[0].End)
	}
}

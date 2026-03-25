package slotbuilder

import (
	"time"

	scheduledomain "github.com/SashaMaltsev/room-booking-service/internal/domain/schedule"
	slotdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/slot"
)

func BuildForDate(schedule scheduledomain.Schedule, date time.Time) []slotdomain.Slot {
	windows := schedule.WindowsForDate(date)
	if len(windows) == 0 {
		return nil
	}

	slots := make([]slotdomain.Slot, 0, len(windows))
	for _, window := range windows {
		slots = append(slots, slotdomain.Slot{
			RoomID:     schedule.RoomID,
			ScheduleID: schedule.ID,
			Start:      window.Start,
			End:        window.End,
		})
	}

	return slots
}

func BuildForHorizon(schedule scheduledomain.Schedule, from time.Time, days int) []slotdomain.Slot {
	if days <= 0 {
		return nil
	}

	start := from.UTC()
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)

	slots := make([]slotdomain.Slot, 0)
	for i := range days {
		date := start.AddDate(0, 0, i)
		slots = append(slots, BuildForDate(schedule, date)...)
	}

	return slots
}

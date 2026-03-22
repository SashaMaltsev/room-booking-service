package slot

import (
	"fmt"
	"strings"
	"time"

	"github.com/SashaMaltsev/room-booking-service/internal/domain/common"
)

type Slot struct {
	ID         string
	RoomID     string
	ScheduleID string
	Start      time.Time
	End        time.Time
	CreatedAt  time.Time
}

func (s *Slot) Normalize() {
	s.RoomID = strings.TrimSpace(s.RoomID)
	s.ScheduleID = strings.TrimSpace(s.ScheduleID)
	s.Start = s.Start.UTC()
	s.End = s.End.UTC()
}

func (s Slot) Validate() error {
	if strings.TrimSpace(s.RoomID) == "" {
		return ErrRoomIDRequired
	}

	if s.Start.IsZero() || s.End.IsZero() {
		return ErrTimeRangeRequired
	}

	start := s.Start.UTC()
	end := s.End.UTC()
	if !start.Before(end) {
		return ErrInvalidTimeRange
	}

	if end.Sub(start) != common.SlotDuration {
		return fmt.Errorf("%w: %s", ErrInvalidSlotDuration, end.Sub(start))
	}

	return nil
}

func (s Slot) Overlaps(other Slot) bool {
	return s.Start.UTC().Before(other.End.UTC()) && other.Start.UTC().Before(s.End.UTC())
}

func (s Slot) IsPast(now time.Time) bool {
	return s.Start.UTC().Before(now.UTC())
}

func (s Slot) Date() time.Time {
	start := s.Start.UTC()
	return time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
}

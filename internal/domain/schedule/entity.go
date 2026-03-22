package schedule

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/SashaMaltsev/room-booking-service/internal/domain/common"
)

type Weekday uint8

const (
	Monday Weekday = 1 + iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

type TimeOfDay struct {
	minutes int
}

type Window struct {
	Start time.Time
	End   time.Time
}

type Schedule struct {
	ID         string
	RoomID     string
	DaysOfWeek []Weekday
	StartTime  TimeOfDay
	EndTime    TimeOfDay
	CreatedAt  time.Time
}

func NewTimeOfDay(hour, minute int) (TimeOfDay, error) {
	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return TimeOfDay{}, fmt.Errorf("%w: %02d:%02d", ErrInvalidTimeOfDay, hour, minute)
	}

	return TimeOfDay{minutes: hour*60 + minute}, nil
}

func ParseTimeOfDay(raw string) (TimeOfDay, error) {
	parsed, err := time.Parse("15:04", raw)
	if err != nil {
		return TimeOfDay{}, fmt.Errorf("%w: %q", ErrInvalidTimeOfDay, raw)
	}

	return NewTimeOfDay(parsed.Hour(), parsed.Minute())
}

func (t TimeOfDay) Hour() int {
	return t.minutes / 60
}

func (t TimeOfDay) Minute() int {
	return t.minutes % 60
}

func (t TimeOfDay) String() string {
	return fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
}

func (t TimeOfDay) Before(other TimeOfDay) bool {
	return t.minutes < other.minutes
}

func (t TimeOfDay) Sub(other TimeOfDay) time.Duration {
	return time.Duration(t.minutes-other.minutes) * time.Minute
}

func (t TimeOfDay) On(date time.Time) time.Time {
	day := date.UTC()
	return time.Date(day.Year(), day.Month(), day.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
}

func (d Weekday) IsValid() bool {
	return d >= Monday && d <= Sunday
}

func (d Weekday) String() string {
	switch d {
	case Monday:
		return "monday"
	case Tuesday:
		return "tuesday"
	case Wednesday:
		return "wednesday"
	case Thursday:
		return "thursday"
	case Friday:
		return "friday"
	case Saturday:
		return "saturday"
	case Sunday:
		return "sunday"
	default:
		return "unknown"
	}
}

func ParseWeekday(value int) (Weekday, error) {
	day := Weekday(value)
	if !day.IsValid() {
		return 0, fmt.Errorf("%w: %d", ErrInvalidDayOfWeek, value)
	}

	return day, nil
}

func ISOWeekday(date time.Time) Weekday {
	switch date.UTC().Weekday() {
	case time.Monday:
		return Monday
	case time.Tuesday:
		return Tuesday
	case time.Wednesday:
		return Wednesday
	case time.Thursday:
		return Thursday
	case time.Friday:
		return Friday
	case time.Saturday:
		return Saturday
	case time.Sunday:
		return Sunday
	default:
		return 0
	}
}

func (s *Schedule) Normalize() {
	s.RoomID = strings.TrimSpace(s.RoomID)
	s.DaysOfWeek = normalizeWeekdays(s.DaysOfWeek)
}

func (s Schedule) Validate() error {
	if strings.TrimSpace(s.RoomID) == "" {
		return ErrRoomIDRequired
	}

	if len(s.DaysOfWeek) == 0 {
		return ErrDaysOfWeekRequired
	}

	seen := make(map[Weekday]struct{}, len(s.DaysOfWeek))
	for _, day := range s.DaysOfWeek {
		if !day.IsValid() {
			return fmt.Errorf("%w: %d", ErrInvalidDayOfWeek, day)
		}

		if _, exists := seen[day]; exists {
			return fmt.Errorf("%w: %d", ErrDuplicateDayOfWeek, day)
		}

		seen[day] = struct{}{}
	}

	if !s.StartTime.Before(s.EndTime) {
		return ErrInvalidTimeRange
	}

	duration := s.EndTime.Sub(s.StartTime)
	if duration < common.SlotDuration || duration%common.SlotDuration != 0 {
		return fmt.Errorf("%w: %s", ErrWindowNotAlignedToSlotDuration, duration)
	}

	return nil
}

func (s Schedule) Duration() time.Duration {
	return s.EndTime.Sub(s.StartTime)
}

func (s Schedule) Covers(day Weekday) bool {
	for _, availableDay := range s.DaysOfWeek {
		if availableDay == day {
			return true
		}
	}

	return false
}

func (s Schedule) MatchesDate(date time.Time) bool {
	return s.Covers(ISOWeekday(date))
}

func (s Schedule) WindowsForDate(date time.Time) []Window {
	if !s.MatchesDate(date) {
		return nil
	}

	start := s.StartTime.On(date)
	end := s.EndTime.On(date)
	count := int(s.Duration() / common.SlotDuration)

	windows := make([]Window, 0, count)
	for current := start; !current.Add(common.SlotDuration).After(end); current = current.Add(common.SlotDuration) {
		next := current.Add(common.SlotDuration)
		windows = append(windows, Window{
			Start: current,
			End:   next,
		})
	}

	return windows
}

func normalizeWeekdays(days []Weekday) []Weekday {
	if len(days) == 0 {
		return nil
	}

	seen := make(map[Weekday]struct{}, len(days))
	normalized := make([]Weekday, 0, len(days))

	for _, day := range days {
		if _, exists := seen[day]; exists {
			continue
		}

		seen[day] = struct{}{}
		normalized = append(normalized, day)
	}

	sort.Slice(normalized, func(i, j int) bool {
		return normalized[i] < normalized[j]
	})

	return normalized
}

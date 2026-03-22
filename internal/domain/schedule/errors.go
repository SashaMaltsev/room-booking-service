package schedule

import "errors"

var (
	ErrScheduleNotFound               = errors.New("schedule not found")
	ErrScheduleExists                 = errors.New("schedule already exists")
	ErrRoomIDRequired                 = errors.New("room id is required")
	ErrDaysOfWeekRequired             = errors.New("days of week are required")
	ErrInvalidDayOfWeek               = errors.New("day of week is invalid")
	ErrDuplicateDayOfWeek             = errors.New("day of week is duplicated")
	ErrInvalidTimeOfDay               = errors.New("time of day is invalid")
	ErrInvalidTimeRange               = errors.New("schedule time range is invalid")
	ErrWindowNotAlignedToSlotDuration = errors.New("schedule window must be aligned to slot duration")
)

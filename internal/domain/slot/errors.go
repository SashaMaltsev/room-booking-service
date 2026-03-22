package slot

import "errors"

var (
	ErrSlotNotFound        = errors.New("slot not found")
	ErrRoomIDRequired      = errors.New("room id is required")
	ErrTimeRangeRequired   = errors.New("slot time range is required")
	ErrInvalidTimeRange    = errors.New("slot time range is invalid")
	ErrInvalidSlotDuration = errors.New("slot duration must be exactly 30 minutes")
	ErrSlotInPast          = errors.New("slot is in the past")
)

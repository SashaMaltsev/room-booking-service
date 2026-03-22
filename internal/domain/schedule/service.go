package schedule

import "context"

type CreateInput struct {
	RoomID     string
	DaysOfWeek []Weekday
	StartTime  TimeOfDay
	EndTime    TimeOfDay
}

type Service interface {
	Create(ctx context.Context, input CreateInput) (Schedule, error)
	GetByRoomID(ctx context.Context, roomID string) (Schedule, error)
}

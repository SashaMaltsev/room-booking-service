package schedule

import "context"

type Repository interface {
	Create(ctx context.Context, schedule Schedule) (Schedule, error)
	GetByRoomID(ctx context.Context, roomID string) (Schedule, error)
	ExistsForRoom(ctx context.Context, roomID string) (bool, error)
}

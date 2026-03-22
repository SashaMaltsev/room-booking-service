package slot

import (
	"context"
	"time"
)

type Repository interface {
	CreateMany(ctx context.Context, slots []Slot) error
	GetByID(ctx context.Context, id string) (Slot, error)
	ListByRoomAndRange(ctx context.Context, roomID string, start, end time.Time) ([]Slot, error)
	ListAvailableByRoomAndDate(ctx context.Context, roomID string, date time.Time) ([]Slot, error)
}

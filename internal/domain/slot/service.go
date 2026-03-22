package slot

import (
	"context"
	"time"
)

type Service interface {
	EnsureGeneratedForDate(ctx context.Context, roomID string, date time.Time) error
	GetByID(ctx context.Context, id string) (Slot, error)
	ListAvailableByRoomAndDate(ctx context.Context, roomID string, date time.Time) ([]Slot, error)
}

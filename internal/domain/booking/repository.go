package booking

import (
	"context"
	"time"

	"github.com/SashaMaltsev/room-booking-service/internal/domain/common"
)

type Repository interface {
	Create(ctx context.Context, booking Booking) (Booking, error)
	GetByID(ctx context.Context, id string) (Booking, error)
	ExistsActiveBySlotID(ctx context.Context, slotID string) (bool, error)
	ListAll(ctx context.Context, page common.PageRequest) ([]Booking, common.Page, error)
	ListUpcomingByUser(ctx context.Context, userID string, now time.Time) ([]Booking, error)
	Update(ctx context.Context, booking Booking) (Booking, error)
}

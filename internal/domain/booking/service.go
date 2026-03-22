package booking

import (
	"context"

	"github.com/SashaMaltsev/room-booking-service/internal/domain/common"
)

type CreateInput struct {
	UserID               string
	SlotID               string
	CreateConferenceLink bool
}

type Service interface {
	Create(ctx context.Context, input CreateInput) (Booking, error)
	Cancel(ctx context.Context, bookingID, userID string) (Booking, error)
	ListAll(ctx context.Context, page common.PageRequest) ([]Booking, common.Page, error)
	ListUpcomingByUser(ctx context.Context, userID string) ([]Booking, error)
}

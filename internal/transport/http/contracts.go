package httptransport

import (
	"context"

	bookingdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/booking"
	commondomain "github.com/SashaMaltsev/room-booking-service/internal/domain/common"
	roomdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/room"
	scheduledomain "github.com/SashaMaltsev/room-booking-service/internal/domain/schedule"
	slotdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/slot"
)

type AuthService interface {
	DummyLogin(ctx context.Context, role commondomain.Role) (string, error)
}

type TokenVerifier interface {
	Parse(token string) (Principal, error)
}

type Dependencies struct {
	Auth          AuthService
	Rooms         roomdomain.Service
	Schedules     scheduledomain.Service
	Slots         slotdomain.Service
	Bookings      bookingdomain.Service
	TokenVerifier TokenVerifier
}

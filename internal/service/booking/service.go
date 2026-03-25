package booking

import (
	"context"

	bookingdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/booking"
	commondomain "github.com/SashaMaltsev/room-booking-service/internal/domain/common"
	slotdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/slot"
)

var _ bookingdomain.Service = (*Service)(nil)

type Service struct {
	bookings bookingdomain.Repository
	slots    slotdomain.Repository
	clock    commondomain.Clock
}

func New(bookings bookingdomain.Repository, slots slotdomain.Repository, clock commondomain.Clock) *Service {
	if clock == nil {
		clock = commondomain.SystemClock{}
	}

	return &Service{
		bookings: bookings,
		slots:    slots,
		clock:    clock,
	}
}

func (s *Service) Create(ctx context.Context, input bookingdomain.CreateInput) (bookingdomain.Booking, error) {
	slot, err := s.slots.GetByID(ctx, input.SlotID)
	if err != nil {
		return bookingdomain.Booking{}, err
	}

	if slot.IsPast(s.clock.Now()) {
		return bookingdomain.Booking{}, slotdomain.ErrSlotInPast
	}

	alreadyBooked, err := s.bookings.ExistsActiveBySlotID(ctx, input.SlotID)
	if err != nil {
		return bookingdomain.Booking{}, err
	}
	if alreadyBooked {
		return bookingdomain.Booking{}, bookingdomain.ErrSlotAlreadyBooked
	}

	entity := bookingdomain.Booking{
		SlotID: input.SlotID,
		UserID: input.UserID,
		Status: bookingdomain.StatusActive,
	}
	entity.Normalize()

	if err := entity.Validate(); err != nil {
		return bookingdomain.Booking{}, err
	}

	return s.bookings.Create(ctx, entity)
}

func (s *Service) Cancel(ctx context.Context, bookingID, userID string) (bookingdomain.Booking, error) {
	entity, err := s.bookings.GetByID(ctx, bookingID)
	if err != nil {
		return bookingdomain.Booking{}, err
	}

	if err := entity.Cancel(userID, s.clock.Now()); err != nil {
		return bookingdomain.Booking{}, err
	}

	return s.bookings.Update(ctx, entity)
}

func (s *Service) ListAll(ctx context.Context, page commondomain.PageRequest) ([]bookingdomain.Booking, commondomain.Page, error) {
	return s.bookings.ListAll(ctx, page)
}

func (s *Service) ListUpcomingByUser(ctx context.Context, userID string) ([]bookingdomain.Booking, error) {
	return s.bookings.ListUpcomingByUser(ctx, userID, s.clock.Now())
}

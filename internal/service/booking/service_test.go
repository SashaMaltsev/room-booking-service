package booking

import (
	"context"
	"errors"
	"testing"
	"time"

	bookingdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/booking"
	commondomain "github.com/SashaMaltsev/room-booking-service/internal/domain/common"
	slotdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/slot"
)

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time {
	return c.now
}

type bookingRepositoryStub struct {
	createFn               func(ctx context.Context, booking bookingdomain.Booking) (bookingdomain.Booking, error)
	getByIDFn              func(ctx context.Context, id string) (bookingdomain.Booking, error)
	existsActiveBySlotIDFn func(ctx context.Context, slotID string) (bool, error)
	listAllFn              func(ctx context.Context, page commondomain.PageRequest) ([]bookingdomain.Booking, commondomain.Page, error)
	listUpcomingByUserFn   func(ctx context.Context, userID string, now time.Time) ([]bookingdomain.Booking, error)
	updateFn               func(ctx context.Context, booking bookingdomain.Booking) (bookingdomain.Booking, error)
}

func (s bookingRepositoryStub) Create(ctx context.Context, booking bookingdomain.Booking) (bookingdomain.Booking, error) {
	if s.createFn != nil {
		return s.createFn(ctx, booking)
	}

	return booking, nil
}

func (s bookingRepositoryStub) GetByID(ctx context.Context, id string) (bookingdomain.Booking, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id)
	}

	return bookingdomain.Booking{}, nil
}

func (s bookingRepositoryStub) ExistsActiveBySlotID(ctx context.Context, slotID string) (bool, error) {
	if s.existsActiveBySlotIDFn != nil {
		return s.existsActiveBySlotIDFn(ctx, slotID)
	}

	return false, nil
}

func (s bookingRepositoryStub) ListAll(ctx context.Context, page commondomain.PageRequest) ([]bookingdomain.Booking, commondomain.Page, error) {
	if s.listAllFn != nil {
		return s.listAllFn(ctx, page)
	}

	return nil, commondomain.Page{}, nil
}

func (s bookingRepositoryStub) ListUpcomingByUser(ctx context.Context, userID string, now time.Time) ([]bookingdomain.Booking, error) {
	if s.listUpcomingByUserFn != nil {
		return s.listUpcomingByUserFn(ctx, userID, now)
	}

	return nil, nil
}

func (s bookingRepositoryStub) Update(ctx context.Context, booking bookingdomain.Booking) (bookingdomain.Booking, error) {
	if s.updateFn != nil {
		return s.updateFn(ctx, booking)
	}

	return booking, nil
}

type slotRepositoryStub struct {
	getByIDFn func(ctx context.Context, id string) (slotdomain.Slot, error)
}

func (s slotRepositoryStub) CreateMany(ctx context.Context, slots []slotdomain.Slot) error {
	return nil
}

func (s slotRepositoryStub) GetByID(ctx context.Context, id string) (slotdomain.Slot, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id)
	}

	return slotdomain.Slot{}, nil
}

func (s slotRepositoryStub) ListByRoomAndRange(ctx context.Context, roomID string, start, end time.Time) ([]slotdomain.Slot, error) {
	return nil, nil
}

func (s slotRepositoryStub) ListAvailableByRoomAndDate(ctx context.Context, roomID string, date time.Time) ([]slotdomain.Slot, error) {
	return nil, nil
}

func TestCreateRejectsPastSlot(t *testing.T) {
	service := New(
		bookingRepositoryStub{},
		slotRepositoryStub{
			getByIDFn: func(_ context.Context, id string) (slotdomain.Slot, error) {
				return slotdomain.Slot{
					ID:    id,
					Start: time.Date(2026, time.March, 24, 8, 0, 0, 0, time.UTC),
					End:   time.Date(2026, time.March, 24, 8, 30, 0, 0, time.UTC),
				}, nil
			},
		},
		fixedClock{now: time.Date(2026, time.March, 24, 9, 0, 0, 0, time.UTC)},
	)

	_, err := service.Create(context.Background(), bookingdomain.CreateInput{
		UserID: "user-1",
		SlotID: "slot-1",
	})
	if !errors.Is(err, slotdomain.ErrSlotInPast) {
		t.Fatalf("expected ErrSlotInPast, got %v", err)
	}
}

func TestCreateRejectsAlreadyBookedSlot(t *testing.T) {
	service := New(
		bookingRepositoryStub{
			existsActiveBySlotIDFn: func(_ context.Context, slotID string) (bool, error) {
				return true, nil
			},
		},
		slotRepositoryStub{
			getByIDFn: func(_ context.Context, id string) (slotdomain.Slot, error) {
				return slotdomain.Slot{
					ID:    id,
					Start: time.Date(2026, time.March, 24, 10, 0, 0, 0, time.UTC),
					End:   time.Date(2026, time.March, 24, 10, 30, 0, 0, time.UTC),
				}, nil
			},
		},
		fixedClock{now: time.Date(2026, time.March, 24, 9, 0, 0, 0, time.UTC)},
	)

	_, err := service.Create(context.Background(), bookingdomain.CreateInput{
		UserID: "user-1",
		SlotID: "slot-1",
	})
	if !errors.Is(err, bookingdomain.ErrSlotAlreadyBooked) {
		t.Fatalf("expected ErrSlotAlreadyBooked, got %v", err)
	}
}

func TestCancelUpdatesBooking(t *testing.T) {
	var updated bookingdomain.Booking

	service := New(
		bookingRepositoryStub{
			getByIDFn: func(_ context.Context, id string) (bookingdomain.Booking, error) {
				return bookingdomain.Booking{
					ID:     id,
					SlotID: "slot-1",
					UserID: "user-1",
					Status: bookingdomain.StatusActive,
				}, nil
			},
			updateFn: func(_ context.Context, booking bookingdomain.Booking) (bookingdomain.Booking, error) {
				updated = booking
				return booking, nil
			},
		},
		slotRepositoryStub{},
		fixedClock{now: time.Date(2026, time.March, 24, 11, 0, 0, 0, time.UTC)},
	)

	cancelled, err := service.Cancel(context.Background(), "booking-1", "user-1")
	if err != nil {
		t.Fatalf("expected cancel without error, got %v", err)
	}

	if updated.Status != bookingdomain.StatusCancelled {
		t.Fatalf("expected updated status cancelled, got %q", updated.Status)
	}

	if cancelled.CancelledAt == nil {
		t.Fatalf("expected cancelledAt to be set")
	}
}

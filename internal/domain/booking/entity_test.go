package booking

import (
	"errors"
	"testing"
	"time"
)

func TestBookingCancelChangesStateAndIsIdempotent(t *testing.T) {
	entity := Booking{
		SlotID: "slot-1",
		UserID: "user-1",
		Status: StatusActive,
	}

	cancelledAt := time.Date(2026, time.March, 23, 10, 0, 0, 0, time.UTC)
	if err := entity.Cancel("user-1", cancelledAt); err != nil {
		t.Fatalf("expected successful cancel, got %v", err)
	}

	if entity.Status != StatusCancelled {
		t.Fatalf("expected status %q, got %q", StatusCancelled, entity.Status)
	}

	if entity.CancelledAt == nil || !entity.CancelledAt.Equal(cancelledAt) {
		t.Fatalf("expected cancelledAt to be %s, got %v", cancelledAt, entity.CancelledAt)
	}

	previousCancelledAt := *entity.CancelledAt
	if err := entity.Cancel("user-1", time.Time{}); err != nil {
		t.Fatalf("expected idempotent second cancel, got %v", err)
	}

	if entity.CancelledAt == nil || !entity.CancelledAt.Equal(previousCancelledAt) {
		t.Fatalf("expected cancelledAt to remain %s, got %v", previousCancelledAt, entity.CancelledAt)
	}
}

func TestBookingCancelRejectsAnotherUser(t *testing.T) {
	entity := Booking{
		SlotID: "slot-1",
		UserID: "owner",
		Status: StatusActive,
	}

	err := entity.Cancel("intruder", time.Date(2026, time.March, 23, 10, 0, 0, 0, time.UTC))
	if !errors.Is(err, ErrCannotCancelAnotherUsersBooking) {
		t.Fatalf("expected ErrCannotCancelAnotherUsersBooking, got %v", err)
	}
}

func TestBookingValidateRequiresCancelledAtForCancelledStatus(t *testing.T) {
	entity := Booking{
		SlotID: "slot-1",
		UserID: "user-1",
		Status: StatusCancelled,
	}

	err := entity.Validate()
	if !errors.Is(err, ErrCancelledAtRequired) {
		t.Fatalf("expected ErrCancelledAtRequired, got %v", err)
	}
}

package booking

import (
	"fmt"
	"strings"
	"time"
)

type Status string

const (
	StatusActive    Status = "active"
	StatusCancelled Status = "cancelled"
)

type Booking struct {
	ID             string
	SlotID         string
	UserID         string
	Status         Status
	ConferenceLink *string
	CreatedAt      time.Time
	CancelledAt    *time.Time
}

func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusCancelled:
		return true
	default:
		return false
	}
}

func (b *Booking) Normalize() {
	b.SlotID = strings.TrimSpace(b.SlotID)
	b.UserID = strings.TrimSpace(b.UserID)
	b.ConferenceLink = normalizeNullableString(b.ConferenceLink)
}

func (b Booking) Validate() error {
	if strings.TrimSpace(b.SlotID) == "" {
		return ErrSlotIDRequired
	}

	if strings.TrimSpace(b.UserID) == "" {
		return ErrUserIDRequired
	}

	if !b.Status.IsValid() {
		return fmt.Errorf("%w: %q", ErrInvalidStatus, b.Status)
	}

	if b.Status == StatusCancelled && b.CancelledAt == nil {
		return ErrCancelledAtRequired
	}

	if b.Status == StatusActive && b.CancelledAt != nil {
		return ErrActiveBookingCannotHaveCancelledAt
	}

	return nil
}

func (b Booking) IsActive() bool {
	return b.Status == StatusActive
}

func (b Booking) IsCancelled() bool {
	return b.Status == StatusCancelled
}

func (b *Booking) Cancel(requestedBy string, cancelledAt time.Time) error {
	if strings.TrimSpace(requestedBy) == "" {
		return ErrUserIDRequired
	}

	if b.UserID != requestedBy {
		return ErrCannotCancelAnotherUsersBooking
	}

	if b.Status == StatusCancelled {
		return nil
	}

	if cancelledAt.IsZero() {
		return ErrCancelledAtRequired
	}

	cancelledAtUTC := cancelledAt.UTC()
	b.Status = StatusCancelled
	b.CancelledAt = &cancelledAtUTC

	return nil
}

func normalizeNullableString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

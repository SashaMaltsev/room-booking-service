package booking

import "errors"

var (
	ErrBookingNotFound                    = errors.New("booking not found")
	ErrSlotIDRequired                     = errors.New("slot id is required")
	ErrUserIDRequired                     = errors.New("user id is required")
	ErrInvalidStatus                      = errors.New("booking status is invalid")
	ErrCancelledAtRequired                = errors.New("cancelled booking must have cancelled_at")
	ErrActiveBookingCannotHaveCancelledAt = errors.New("active booking cannot have cancelled_at")
	ErrCannotCancelAnotherUsersBooking    = errors.New("cannot cancel another user's booking")
	ErrSlotAlreadyBooked                  = errors.New("slot is already booked")
)

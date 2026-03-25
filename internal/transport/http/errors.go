package httptransport

import (
	"errors"
	"net/http"

	bookingdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/booking"
	commondomain "github.com/SashaMaltsev/room-booking-service/internal/domain/common"
	roomdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/room"
	scheduledomain "github.com/SashaMaltsev/room-booking-service/internal/domain/schedule"
	slotdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/slot"
)

type errorBody struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, roomdomain.ErrRoomNotFound):
		writeJSON(w, http.StatusNotFound, errorBody{Error: apiError{Code: "ROOM_NOT_FOUND", Message: err.Error()}})
	case errors.Is(err, slotdomain.ErrSlotNotFound):
		writeJSON(w, http.StatusNotFound, errorBody{Error: apiError{Code: "SLOT_NOT_FOUND", Message: err.Error()}})
	case errors.Is(err, bookingdomain.ErrBookingNotFound):
		writeJSON(w, http.StatusNotFound, errorBody{Error: apiError{Code: "BOOKING_NOT_FOUND", Message: err.Error()}})
	case errors.Is(err, scheduledomain.ErrScheduleExists):
		writeJSON(w, http.StatusConflict, errorBody{Error: apiError{Code: "SCHEDULE_EXISTS", Message: err.Error()}})
	case errors.Is(err, bookingdomain.ErrSlotAlreadyBooked):
		writeJSON(w, http.StatusConflict, errorBody{Error: apiError{Code: "SLOT_ALREADY_BOOKED", Message: err.Error()}})
	case errors.Is(err, bookingdomain.ErrCannotCancelAnotherUsersBooking):
		writeJSON(w, http.StatusForbidden, errorBody{Error: apiError{Code: "FORBIDDEN", Message: err.Error()}})
	case errors.Is(err, commondomain.ErrInvalidRole),
		errors.Is(err, roomdomain.ErrNameRequired),
		errors.Is(err, roomdomain.ErrCapacityMustBePositive),
		errors.Is(err, scheduledomain.ErrRoomIDRequired),
		errors.Is(err, scheduledomain.ErrDaysOfWeekRequired),
		errors.Is(err, scheduledomain.ErrInvalidDayOfWeek),
		errors.Is(err, scheduledomain.ErrDuplicateDayOfWeek),
		errors.Is(err, scheduledomain.ErrInvalidTimeOfDay),
		errors.Is(err, scheduledomain.ErrInvalidTimeRange),
		errors.Is(err, scheduledomain.ErrWindowNotAlignedToSlotDuration),
		errors.Is(err, slotdomain.ErrRoomIDRequired),
		errors.Is(err, slotdomain.ErrTimeRangeRequired),
		errors.Is(err, slotdomain.ErrInvalidTimeRange),
		errors.Is(err, slotdomain.ErrInvalidSlotDuration),
		errors.Is(err, slotdomain.ErrSlotInPast),
		errors.Is(err, bookingdomain.ErrSlotIDRequired),
		errors.Is(err, bookingdomain.ErrUserIDRequired),
		errors.Is(err, bookingdomain.ErrInvalidStatus),
		errors.Is(err, bookingdomain.ErrCancelledAtRequired),
		errors.Is(err, bookingdomain.ErrActiveBookingCannotHaveCancelledAt),
		errors.Is(err, commondomain.ErrInvalidPage),
		errors.Is(err, commondomain.ErrInvalidPageSize):
		writeJSON(w, http.StatusBadRequest, errorBody{Error: apiError{Code: "INVALID_REQUEST", Message: err.Error()}})
	default:
		writeJSON(w, http.StatusInternalServerError, errorBody{Error: apiError{Code: "INTERNAL_ERROR", Message: "internal server error"}})
	}
}

func writeUnauthorized(w http.ResponseWriter) {
	writeJSON(w, http.StatusUnauthorized, errorBody{Error: apiError{Code: "UNAUTHORIZED", Message: "unauthorized"}})
}

func writeForbidden(w http.ResponseWriter) {
	writeJSON(w, http.StatusForbidden, errorBody{Error: apiError{Code: "FORBIDDEN", Message: "forbidden"}})
}

func writeInvalidRequest(w http.ResponseWriter, message string) {
	writeJSON(w, http.StatusBadRequest, errorBody{Error: apiError{Code: "INVALID_REQUEST", Message: message}})
}

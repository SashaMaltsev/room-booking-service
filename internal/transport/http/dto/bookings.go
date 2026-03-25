package dto

import (
	"time"

	bookingdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/booking"
)

type CreateBookingRequest struct {
	SlotID               string `json:"slotId"`
	CreateConferenceLink bool   `json:"createConferenceLink"`
}

type BookingResponse struct {
	ID             string     `json:"id"`
	SlotID         string     `json:"slotId"`
	UserID         string     `json:"userId"`
	Status         string     `json:"status"`
	ConferenceLink *string    `json:"conferenceLink"`
	CreatedAt      *time.Time `json:"createdAt"`
}

type BookingEnvelope struct {
	Booking BookingResponse `json:"booking"`
}

type ListBookingsResponse struct {
	Bookings   []BookingResponse   `json:"bookings"`
	Pagination *PaginationResponse `json:"pagination,omitempty"`
}

type PaginationResponse struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

func NewBookingResponse(entity bookingdomain.Booking) BookingResponse {
	return BookingResponse{
		ID:             entity.ID,
		SlotID:         entity.SlotID,
		UserID:         entity.UserID,
		Status:         string(entity.Status),
		ConferenceLink: entity.ConferenceLink,
		CreatedAt:      nullableTime(entity.CreatedAt),
	}
}

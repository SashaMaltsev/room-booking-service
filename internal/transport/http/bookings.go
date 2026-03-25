package httptransport

import (
	"net/http"

	bookingdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/booking"
	"github.com/SashaMaltsev/room-booking-service/internal/transport/http/dto"
)

func (a *API) handleCreateBooking(w http.ResponseWriter, r *http.Request) {
	if a.bookings == nil {
		writeDomainError(w, nil)
		return
	}

	principal, ok := PrincipalFromContext(r.Context())
	if !ok {
		writeUnauthorized(w)
		return
	}

	var request dto.CreateBookingRequest
	if err := decodeJSON(r, &request); err != nil {
		writeInvalidRequest(w, "invalid request body")
		return
	}

	booking, err := a.bookings.Create(r.Context(), bookingdomain.CreateInput{
		UserID:               principal.UserID,
		SlotID:               request.SlotID,
		CreateConferenceLink: request.CreateConferenceLink,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, dto.BookingEnvelope{
		Booking: dto.NewBookingResponse(booking),
	})
}

func (a *API) handleListBookings(w http.ResponseWriter, r *http.Request) {
	if a.bookings == nil {
		writeDomainError(w, nil)
		return
	}

	page, err := parsePagination(r)
	if err != nil {
		writeInvalidRequest(w, "invalid pagination parameters")
		return
	}

	bookings, pagination, err := a.bookings.ListAll(r.Context(), page)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	response := dto.ListBookingsResponse{
		Bookings: make([]dto.BookingResponse, 0, len(bookings)),
		Pagination: &dto.PaginationResponse{
			Page:     pagination.Page,
			PageSize: pagination.PageSize,
			Total:    pagination.Total,
		},
	}
	for _, entity := range bookings {
		response.Bookings = append(response.Bookings, dto.NewBookingResponse(entity))
	}

	writeJSON(w, http.StatusOK, response)
}

func (a *API) handleListMyBookings(w http.ResponseWriter, r *http.Request) {
	if a.bookings == nil {
		writeDomainError(w, nil)
		return
	}

	principal, ok := PrincipalFromContext(r.Context())
	if !ok {
		writeUnauthorized(w)
		return
	}

	bookings, err := a.bookings.ListUpcomingByUser(r.Context(), principal.UserID)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	response := dto.ListBookingsResponse{
		Bookings: make([]dto.BookingResponse, 0, len(bookings)),
	}
	for _, entity := range bookings {
		response.Bookings = append(response.Bookings, dto.NewBookingResponse(entity))
	}

	writeJSON(w, http.StatusOK, response)
}

func (a *API) handleCancelBooking(w http.ResponseWriter, r *http.Request) {
	if a.bookings == nil {
		writeDomainError(w, nil)
		return
	}

	principal, ok := PrincipalFromContext(r.Context())
	if !ok {
		writeUnauthorized(w)
		return
	}

	bookingID := r.PathValue("bookingId")
	if bookingID == "" {
		writeInvalidRequest(w, "bookingId is required")
		return
	}

	booking, err := a.bookings.Cancel(r.Context(), bookingID, principal.UserID)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.BookingEnvelope{
		Booking: dto.NewBookingResponse(booking),
	})
}

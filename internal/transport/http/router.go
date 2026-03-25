package httptransport

import (
	"net/http"

	bookingdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/booking"
	roomdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/room"
	scheduledomain "github.com/SashaMaltsev/room-booking-service/internal/domain/schedule"
	slotdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/slot"
)

type API struct {
	auth      AuthService
	rooms     roomdomain.Service
	schedules scheduledomain.Service
	slots     slotdomain.Service
	bookings  bookingdomain.Service
	tokens    TokenVerifier
}

func NewHandler(deps Dependencies) http.Handler {
	api := &API{
		auth:      deps.Auth,
		rooms:     deps.Rooms,
		schedules: deps.Schedules,
		slots:     deps.Slots,
		bookings:  deps.Bookings,
		tokens:    deps.TokenVerifier,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", api.handleRoot)
	mux.HandleFunc("GET /_info", api.handleInfo)
	mux.HandleFunc("POST /dummyLogin", api.handleDummyLogin)

	mux.Handle("GET /rooms/list", chain(http.HandlerFunc(api.handleListRooms), api.withAuth))
	mux.Handle("POST /rooms/create", chain(http.HandlerFunc(api.handleCreateRoom), api.withAuth, api.requireRole("admin")))
	mux.Handle("POST /rooms/{roomId}/schedule/create", chain(http.HandlerFunc(api.handleCreateSchedule), api.withAuth, api.requireRole("admin")))
	mux.Handle("GET /rooms/{roomId}/slots/list", chain(http.HandlerFunc(api.handleListSlots), api.withAuth))

	mux.Handle("POST /bookings/create", chain(http.HandlerFunc(api.handleCreateBooking), api.withAuth, api.requireRole("user")))
	mux.Handle("GET /bookings/list", chain(http.HandlerFunc(api.handleListBookings), api.withAuth, api.requireRole("admin")))
	mux.Handle("GET /bookings/my", chain(http.HandlerFunc(api.handleListMyBookings), api.withAuth, api.requireRole("user")))
	mux.Handle("POST /bookings/{bookingId}/cancel", chain(http.HandlerFunc(api.handleCancelBooking), api.withAuth, api.requireRole("user")))

	return mux
}

package httptransport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	bookingdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/booking"
	commondomain "github.com/SashaMaltsev/room-booking-service/internal/domain/common"
	roomdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/room"
	scheduledomain "github.com/SashaMaltsev/room-booking-service/internal/domain/schedule"
	slotdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/slot"
)

type authServiceStub struct {
	dummyLoginFn func(ctx context.Context, role commondomain.Role, demoUser string) (string, error)
}

func (s authServiceStub) DummyLogin(ctx context.Context, role commondomain.Role, demoUser string) (string, error) {
	if s.dummyLoginFn != nil {
		return s.dummyLoginFn(ctx, role, demoUser)
	}

	return "token", nil
}

type tokenVerifierStub struct {
	parseFn func(token string) (Principal, error)
}

func (s tokenVerifierStub) Parse(token string) (Principal, error) {
	if s.parseFn != nil {
		return s.parseFn(token)
	}

	return Principal{}, nil
}

type roomServiceStub struct {
	createFn  func(ctx context.Context, input roomdomain.CreateInput) (roomdomain.Room, error)
	getByIDFn func(ctx context.Context, id string) (roomdomain.Room, error)
	listFn    func(ctx context.Context) ([]roomdomain.Room, error)
}

func (s roomServiceStub) Create(ctx context.Context, input roomdomain.CreateInput) (roomdomain.Room, error) {
	if s.createFn != nil {
		return s.createFn(ctx, input)
	}

	return roomdomain.Room{}, nil
}

func (s roomServiceStub) GetByID(ctx context.Context, id string) (roomdomain.Room, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id)
	}

	return roomdomain.Room{}, nil
}

func (s roomServiceStub) List(ctx context.Context) ([]roomdomain.Room, error) {
	if s.listFn != nil {
		return s.listFn(ctx)
	}

	return nil, nil
}

type scheduleServiceStub struct {
	createFn      func(ctx context.Context, input scheduledomain.CreateInput) (scheduledomain.Schedule, error)
	getByRoomIDFn func(ctx context.Context, roomID string) (scheduledomain.Schedule, error)
}

func (s scheduleServiceStub) Create(ctx context.Context, input scheduledomain.CreateInput) (scheduledomain.Schedule, error) {
	if s.createFn != nil {
		return s.createFn(ctx, input)
	}

	return scheduledomain.Schedule{}, nil
}

func (s scheduleServiceStub) GetByRoomID(ctx context.Context, roomID string) (scheduledomain.Schedule, error) {
	if s.getByRoomIDFn != nil {
		return s.getByRoomIDFn(ctx, roomID)
	}

	return scheduledomain.Schedule{}, nil
}

type slotServiceStub struct {
	ensureGeneratedFn     func(ctx context.Context, roomID string, date time.Time) error
	getByIDFn             func(ctx context.Context, id string) (slotdomain.Slot, error)
	listAvailableByDateFn func(ctx context.Context, roomID string, date time.Time) ([]slotdomain.Slot, error)
}

func (s slotServiceStub) EnsureGeneratedForDate(ctx context.Context, roomID string, date time.Time) error {
	if s.ensureGeneratedFn != nil {
		return s.ensureGeneratedFn(ctx, roomID, date)
	}

	return nil
}

func (s slotServiceStub) GetByID(ctx context.Context, id string) (slotdomain.Slot, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id)
	}

	return slotdomain.Slot{}, nil
}

func (s slotServiceStub) ListAvailableByRoomAndDate(ctx context.Context, roomID string, date time.Time) ([]slotdomain.Slot, error) {
	if s.listAvailableByDateFn != nil {
		return s.listAvailableByDateFn(ctx, roomID, date)
	}

	return nil, nil
}

type bookingServiceStub struct {
	createFn             func(ctx context.Context, input bookingdomain.CreateInput) (bookingdomain.Booking, error)
	cancelFn             func(ctx context.Context, bookingID, userID string) (bookingdomain.Booking, error)
	listAllFn            func(ctx context.Context, page commondomain.PageRequest) ([]bookingdomain.Booking, commondomain.Page, error)
	listUpcomingByUserFn func(ctx context.Context, userID string) ([]bookingdomain.Booking, error)
}

func (s bookingServiceStub) Create(ctx context.Context, input bookingdomain.CreateInput) (bookingdomain.Booking, error) {
	if s.createFn != nil {
		return s.createFn(ctx, input)
	}

	return bookingdomain.Booking{}, nil
}

func (s bookingServiceStub) Cancel(ctx context.Context, bookingID, userID string) (bookingdomain.Booking, error) {
	if s.cancelFn != nil {
		return s.cancelFn(ctx, bookingID, userID)
	}

	return bookingdomain.Booking{}, nil
}

func (s bookingServiceStub) ListAll(ctx context.Context, page commondomain.PageRequest) ([]bookingdomain.Booking, commondomain.Page, error) {
	if s.listAllFn != nil {
		return s.listAllFn(ctx, page)
	}

	return nil, commondomain.Page{}, nil
}

func (s bookingServiceStub) ListUpcomingByUser(ctx context.Context, userID string) ([]bookingdomain.Booking, error) {
	if s.listUpcomingByUserFn != nil {
		return s.listUpcomingByUserFn(ctx, userID)
	}

	return nil, nil
}

func TestDummyLoginReturnsToken(t *testing.T) {
	handler := NewHandler(Dependencies{
		Auth: authServiceStub{
			dummyLoginFn: func(_ context.Context, role commondomain.Role, demoUser string) (string, error) {
				if role != commondomain.RoleAdmin {
					t.Fatalf("expected admin role, got %q", role)
				}
				if demoUser != "user2" {
					t.Fatalf("expected demo user to be forwarded, got %q", demoUser)
				}

				return "signed-token", nil
			},
		},
	})

	request := httptest.NewRequest(http.MethodPost, "/dummyLogin", strings.NewReader(`{"role":"admin","user":"user2"}`))
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("expected valid json response, got %v", err)
	}

	if response["token"] != "signed-token" {
		t.Fatalf("expected token response, got %v", response)
	}
}

func TestCreateRoomRequiresAdminRole(t *testing.T) {
	handler := NewHandler(Dependencies{
		Rooms: roomServiceStub{},
		TokenVerifier: tokenVerifierStub{
			parseFn: func(token string) (Principal, error) {
				return Principal{UserID: "user-1", Role: commondomain.RoleUser}, nil
			},
		},
	})

	request := httptest.NewRequest(http.MethodPost, "/rooms/create", strings.NewReader(`{"name":"Room A"}`))
	request.Header.Set("Authorization", "Bearer token")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", recorder.Code)
	}
}

func TestCreateRoomReturnsCreatedEntity(t *testing.T) {
	handler := NewHandler(Dependencies{
		Rooms: roomServiceStub{
			createFn: func(_ context.Context, input roomdomain.CreateInput) (roomdomain.Room, error) {
				if input.Name != "Room A" {
					t.Fatalf("expected room name to be forwarded, got %q", input.Name)
				}

				return roomdomain.Room{ID: "room-1", Name: input.Name}, nil
			},
		},
		TokenVerifier: tokenVerifierStub{
			parseFn: func(token string) (Principal, error) {
				return Principal{UserID: "admin-1", Role: commondomain.RoleAdmin}, nil
			},
		},
	})

	request := httptest.NewRequest(http.MethodPost, "/rooms/create", strings.NewReader(`{"name":"Room A"}`))
	request.Header.Set("Authorization", "Bearer token")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", recorder.Code)
	}
}

func TestListSlotsRequiresDateQuery(t *testing.T) {
	handler := NewHandler(Dependencies{
		Slots: slotServiceStub{},
		TokenVerifier: tokenVerifierStub{
			parseFn: func(token string) (Principal, error) {
				return Principal{UserID: "user-1", Role: commondomain.RoleUser}, nil
			},
		},
	})

	request := httptest.NewRequest(http.MethodGet, "/rooms/room-1/slots/list", nil)
	request.Header.Set("Authorization", "Bearer token")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestGetSlotReturnsSlot(t *testing.T) {
	start := time.Date(2026, 6, 8, 9, 0, 0, 0, time.UTC)
	handler := NewHandler(Dependencies{
		Slots: slotServiceStub{
			getByIDFn: func(_ context.Context, id string) (slotdomain.Slot, error) {
				if id != "slot-1" {
					t.Fatalf("expected slot id from path, got %q", id)
				}

				return slotdomain.Slot{
					ID:     id,
					RoomID: "room-1",
					Start:  start,
					End:    start.Add(30 * time.Minute),
				}, nil
			},
		},
		TokenVerifier: tokenVerifierStub{
			parseFn: func(token string) (Principal, error) {
				return Principal{UserID: "user-1", Role: commondomain.RoleUser}, nil
			},
		},
	})

	request := httptest.NewRequest(http.MethodGet, "/slots/slot-1", nil)
	request.Header.Set("Authorization", "Bearer token")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	var response struct {
		Slot struct {
			ID     string `json:"id"`
			RoomID string `json:"roomId"`
		} `json:"slot"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("expected valid json response, got %v", err)
	}

	if response.Slot.ID != "slot-1" || response.Slot.RoomID != "room-1" {
		t.Fatalf("expected slot envelope, got %+v", response)
	}
}

func TestCreateBookingUsesPrincipalUserID(t *testing.T) {
	handler := NewHandler(Dependencies{
		Bookings: bookingServiceStub{
			createFn: func(_ context.Context, input bookingdomain.CreateInput) (bookingdomain.Booking, error) {
				if input.UserID != "user-from-token" {
					t.Fatalf("expected user id from token, got %q", input.UserID)
				}

				return bookingdomain.Booking{
					ID:     "booking-1",
					SlotID: input.SlotID,
					UserID: input.UserID,
					Status: bookingdomain.StatusActive,
				}, nil
			},
		},
		TokenVerifier: tokenVerifierStub{
			parseFn: func(token string) (Principal, error) {
				return Principal{UserID: "user-from-token", Role: commondomain.RoleUser}, nil
			},
		},
	})

	request := httptest.NewRequest(http.MethodPost, "/bookings/create", strings.NewReader(`{"slotId":"slot-1","createConferenceLink":true}`))
	request.Header.Set("Authorization", "Bearer token")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", recorder.Code)
	}
}

func TestCreateScheduleMapsConflictError(t *testing.T) {
	handler := NewHandler(Dependencies{
		Schedules: scheduleServiceStub{
			createFn: func(_ context.Context, input scheduledomain.CreateInput) (scheduledomain.Schedule, error) {
				return scheduledomain.Schedule{}, scheduledomain.ErrScheduleExists
			},
		},
		TokenVerifier: tokenVerifierStub{
			parseFn: func(token string) (Principal, error) {
				return Principal{UserID: "admin-1", Role: commondomain.RoleAdmin}, nil
			},
		},
	})

	request := httptest.NewRequest(http.MethodPost, "/rooms/room-1/schedule/create", strings.NewReader(`{"roomId":"room-1","daysOfWeek":[1],"startTime":"09:00","endTime":"10:00"}`))
	request.Header.Set("Authorization", "Bearer token")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", recorder.Code)
	}
}

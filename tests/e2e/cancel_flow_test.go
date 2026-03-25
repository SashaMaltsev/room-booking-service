package e2e

import (
	"net/http"
	"testing"
)

func TestCancelFlow(t *testing.T) {
	baseURL := requireE2E(t)
	client := &http.Client{}

	adminToken := mustLogin(t, client, baseURL, "admin")
	userToken := mustLogin(t, client, baseURL, "user")

	createRoomResp := requestJSON(t, client, http.MethodPost, baseURL+"/rooms/create", adminToken, map[string]any{
		"name": uniqueRoomName("cancel-room"),
	})
	if createRoomResp.StatusCode != http.StatusCreated {
		t.Fatalf("create room expected 201, got %d", createRoomResp.StatusCode)
	}

	var createRoomBody struct {
		Room struct {
			ID string `json:"id"`
		} `json:"room"`
	}
	createRoomBody = decodeBody[struct {
		Room struct {
			ID string `json:"id"`
		} `json:"room"`
	}](t, createRoomResp)

	scheduleResp := requestJSON(t, client, http.MethodPost, baseURL+"/rooms/"+createRoomBody.Room.ID+"/schedule/create", adminToken, map[string]any{
		"roomId":     createRoomBody.Room.ID,
		"daysOfWeek": []int{1, 2, 3, 4, 5, 6, 7},
		"startTime":  "09:00",
		"endTime":    "11:00",
	})
	if scheduleResp.StatusCode != http.StatusCreated {
		t.Fatalf("create schedule expected 201, got %d", scheduleResp.StatusCode)
	}
	_ = decodeBody[map[string]any](t, scheduleResp)

	slotsResp := requestJSON(t, client, http.MethodGet, baseURL+"/rooms/"+createRoomBody.Room.ID+"/slots/list?date="+tomorrowDateUTC(), userToken, nil)
	if slotsResp.StatusCode != http.StatusOK {
		t.Fatalf("list slots expected 200, got %d", slotsResp.StatusCode)
	}

	var slotsBody struct {
		Slots []struct {
			ID string `json:"id"`
		} `json:"slots"`
	}
	slotsBody = decodeBody[struct {
		Slots []struct {
			ID string `json:"id"`
		} `json:"slots"`
	}](t, slotsResp)

	if len(slotsBody.Slots) == 0 {
		t.Fatalf("expected at least one available slot")
	}

	createBookingResp := requestJSON(t, client, http.MethodPost, baseURL+"/bookings/create", userToken, map[string]any{
		"slotId": slotsBody.Slots[0].ID,
	})
	if createBookingResp.StatusCode != http.StatusCreated {
		t.Fatalf("create booking expected 201, got %d", createBookingResp.StatusCode)
	}

	var createBookingBody struct {
		Booking struct {
			ID string `json:"id"`
		} `json:"booking"`
	}
	createBookingBody = decodeBody[struct {
		Booking struct {
			ID string `json:"id"`
		} `json:"booking"`
	}](t, createBookingResp)

	cancelResp := requestJSON(t, client, http.MethodPost, baseURL+"/bookings/"+createBookingBody.Booking.ID+"/cancel", userToken, nil)
	if cancelResp.StatusCode != http.StatusOK {
		t.Fatalf("cancel booking expected 200, got %d", cancelResp.StatusCode)
	}

	var cancelBody struct {
		Booking struct {
			Status string `json:"status"`
		} `json:"booking"`
	}
	cancelBody = decodeBody[struct {
		Booking struct {
			Status string `json:"status"`
		} `json:"booking"`
	}](t, cancelResp)

	if cancelBody.Booking.Status != "cancelled" {
		t.Fatalf("expected cancelled status, got %q", cancelBody.Booking.Status)
	}
}

package httptransport

import (
	"net/http"

	"github.com/SashaMaltsev/room-booking-service/internal/transport/http/dto"
)

func (a *API) handleListSlots(w http.ResponseWriter, r *http.Request) {
	if a.slots == nil {
		writeDomainError(w, nil)
		return
	}

	roomID := r.PathValue("roomId")
	if roomID == "" {
		writeInvalidRequest(w, "roomId is required")
		return
	}

	dateRaw := r.URL.Query().Get("date")
	if dateRaw == "" {
		writeInvalidRequest(w, "date is required")
		return
	}

	date, err := parseDate(dateRaw)
	if err != nil {
		writeInvalidRequest(w, "date must be in YYYY-MM-DD format")
		return
	}

	slots, err := a.slots.ListAvailableByRoomAndDate(r.Context(), roomID, date)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	response := dto.ListSlotsResponse{
		Slots: make([]dto.SlotResponse, 0, len(slots)),
	}
	for _, entity := range slots {
		response.Slots = append(response.Slots, dto.NewSlotResponse(entity))
	}

	writeJSON(w, http.StatusOK, response)
}

func (a *API) handleGetSlot(w http.ResponseWriter, r *http.Request) {
	if a.slots == nil {
		writeDomainError(w, nil)
		return
	}

	slotID := r.PathValue("slotId")
	if slotID == "" {
		writeInvalidRequest(w, "slotId is required")
		return
	}

	slot, err := a.slots.GetByID(r.Context(), slotID)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.SlotEnvelope{
		Slot: dto.NewSlotResponse(slot),
	})
}

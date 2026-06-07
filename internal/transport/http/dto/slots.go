package dto

import (
	"time"

	slotdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/slot"
)

type SlotResponse struct {
	ID     string    `json:"id"`
	RoomID string    `json:"roomId"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
}

type ListSlotsResponse struct {
	Slots []SlotResponse `json:"slots"`
}

type SlotEnvelope struct {
	Slot SlotResponse `json:"slot"`
}

func NewSlotResponse(entity slotdomain.Slot) SlotResponse {
	return SlotResponse{
		ID:     entity.ID,
		RoomID: entity.RoomID,
		Start:  entity.Start.UTC(),
		End:    entity.End.UTC(),
	}
}

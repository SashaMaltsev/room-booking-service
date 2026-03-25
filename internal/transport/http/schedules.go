package httptransport

import (
	"net/http"

	scheduledomain "github.com/SashaMaltsev/room-booking-service/internal/domain/schedule"
	"github.com/SashaMaltsev/room-booking-service/internal/transport/http/dto"
)

func (a *API) handleCreateSchedule(w http.ResponseWriter, r *http.Request) {
	if a.schedules == nil {
		writeDomainError(w, nil)
		return
	}

	roomID := r.PathValue("roomId")
	if roomID == "" {
		writeInvalidRequest(w, "roomId is required")
		return
	}

	var request dto.CreateScheduleRequest
	if err := decodeJSON(r, &request); err != nil {
		writeInvalidRequest(w, "invalid request body")
		return
	}

	if request.RoomID != "" && request.RoomID != roomID {
		writeInvalidRequest(w, "roomId in path and body must match")
		return
	}

	daysOfWeek, err := parseWeekdays(request.DaysOfWeek)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	startTime, err := scheduledomain.ParseTimeOfDay(request.StartTime)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	endTime, err := scheduledomain.ParseTimeOfDay(request.EndTime)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	schedule, err := a.schedules.Create(r.Context(), scheduledomain.CreateInput{
		RoomID:     roomID,
		DaysOfWeek: daysOfWeek,
		StartTime:  startTime,
		EndTime:    endTime,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, dto.CreateScheduleResponse{
		Schedule: dto.NewScheduleResponse(schedule),
	})
}

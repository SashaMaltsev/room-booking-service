package httptransport

import (
	"net/http"

	roomdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/room"
	"github.com/SashaMaltsev/room-booking-service/internal/transport/http/dto"
)

func (a *API) handleListRooms(w http.ResponseWriter, r *http.Request) {
	if a.rooms == nil {
		writeDomainError(w, nil)
		return
	}

	rooms, err := a.rooms.List(r.Context())
	if err != nil {
		writeDomainError(w, err)
		return
	}

	response := dto.ListRoomsResponse{
		Rooms: make([]dto.RoomResponse, 0, len(rooms)),
	}
	for _, entity := range rooms {
		response.Rooms = append(response.Rooms, dto.NewRoomResponse(entity))
	}

	writeJSON(w, http.StatusOK, response)
}

func (a *API) handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	if a.rooms == nil {
		writeDomainError(w, nil)
		return
	}

	var request dto.CreateRoomRequest
	if err := decodeJSON(r, &request); err != nil {
		writeInvalidRequest(w, "invalid request body")
		return
	}

	room, err := a.rooms.Create(r.Context(), roomdomain.CreateInput{
		Name:        request.Name,
		Description: request.Description,
		Capacity:    request.Capacity,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, dto.CreateRoomResponse{
		Room: dto.NewRoomResponse(room),
	})
}

package dto

import (
	"time"

	roomdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/room"
)

type CreateRoomRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Capacity    *int    `json:"capacity"`
}

type RoomResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Capacity    *int       `json:"capacity"`
	CreatedAt   *time.Time `json:"createdAt"`
}

type CreateRoomResponse struct {
	Room RoomResponse `json:"room"`
}

type ListRoomsResponse struct {
	Rooms []RoomResponse `json:"rooms"`
}

func NewRoomResponse(entity roomdomain.Room) RoomResponse {
	return RoomResponse{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		Capacity:    entity.Capacity,
		CreatedAt:   nullableTime(entity.CreatedAt),
	}
}

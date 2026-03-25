package dto

import (
	"time"

	scheduledomain "github.com/SashaMaltsev/room-booking-service/internal/domain/schedule"
)

type CreateScheduleRequest struct {
	RoomID     string `json:"roomId"`
	DaysOfWeek []int  `json:"daysOfWeek"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

type ScheduleResponse struct {
	ID         string     `json:"id,omitempty"`
	RoomID     string     `json:"roomId"`
	DaysOfWeek []int      `json:"daysOfWeek"`
	StartTime  string     `json:"startTime"`
	EndTime    string     `json:"endTime"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
}

type CreateScheduleResponse struct {
	Schedule ScheduleResponse `json:"schedule"`
}

func NewScheduleResponse(entity scheduledomain.Schedule) ScheduleResponse {
	days := make([]int, 0, len(entity.DaysOfWeek))
	for _, day := range entity.DaysOfWeek {
		days = append(days, int(day))
	}

	return ScheduleResponse{
		ID:         entity.ID,
		RoomID:     entity.RoomID,
		DaysOfWeek: days,
		StartTime:  entity.StartTime.String(),
		EndTime:    entity.EndTime.String(),
		CreatedAt:  nullableTime(entity.CreatedAt),
	}
}

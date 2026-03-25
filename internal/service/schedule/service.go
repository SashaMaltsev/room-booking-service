package schedule

import (
	"context"

	commondomain "github.com/SashaMaltsev/room-booking-service/internal/domain/common"
	roomdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/room"
	scheduledomain "github.com/SashaMaltsev/room-booking-service/internal/domain/schedule"
	slotdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/slot"
	"github.com/SashaMaltsev/room-booking-service/internal/service/slotbuilder"
)

const defaultGenerationHorizonDays = 30

var _ scheduledomain.Service = (*Service)(nil)

type Service struct {
	schedules             scheduledomain.Repository
	rooms                 roomdomain.Repository
	slots                 slotdomain.Repository
	clock                 commondomain.Clock
	generationHorizonDays int
}

func New(
	schedules scheduledomain.Repository,
	rooms roomdomain.Repository,
	slots slotdomain.Repository,
	clock commondomain.Clock,
	generationHorizonDays int,
) *Service {
	if clock == nil {
		clock = commondomain.SystemClock{}
	}

	if generationHorizonDays <= 0 {
		generationHorizonDays = defaultGenerationHorizonDays
	}

	return &Service{
		schedules:             schedules,
		rooms:                 rooms,
		slots:                 slots,
		clock:                 clock,
		generationHorizonDays: generationHorizonDays,
	}
}

func (s *Service) Create(ctx context.Context, input scheduledomain.CreateInput) (scheduledomain.Schedule, error) {
	if _, err := s.rooms.GetByID(ctx, input.RoomID); err != nil {
		return scheduledomain.Schedule{}, err
	}

	exists, err := s.schedules.ExistsForRoom(ctx, input.RoomID)
	if err != nil {
		return scheduledomain.Schedule{}, err
	}
	if exists {
		return scheduledomain.Schedule{}, scheduledomain.ErrScheduleExists
	}

	entity := scheduledomain.Schedule{
		RoomID:     input.RoomID,
		DaysOfWeek: input.DaysOfWeek,
		StartTime:  input.StartTime,
		EndTime:    input.EndTime,
	}
	entity.Normalize()

	if err := entity.Validate(); err != nil {
		return scheduledomain.Schedule{}, err
	}

	created, err := s.schedules.Create(ctx, entity)
	if err != nil {
		return scheduledomain.Schedule{}, err
	}

	if s.slots != nil {
		slots := slotbuilder.BuildForHorizon(created, s.clock.Now(), s.generationHorizonDays)
		if err := s.slots.CreateMany(ctx, slots); err != nil {
			return scheduledomain.Schedule{}, err
		}
	}

	return created, nil
}

func (s *Service) GetByRoomID(ctx context.Context, roomID string) (scheduledomain.Schedule, error) {
	return s.schedules.GetByRoomID(ctx, roomID)
}

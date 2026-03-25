package slot

import (
	"context"
	"errors"
	"time"

	commondomain "github.com/SashaMaltsev/room-booking-service/internal/domain/common"
	roomdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/room"
	scheduledomain "github.com/SashaMaltsev/room-booking-service/internal/domain/schedule"
	slotdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/slot"
	"github.com/SashaMaltsev/room-booking-service/internal/service/slotbuilder"
)

var _ slotdomain.Service = (*Service)(nil)

type Service struct {
	rooms     roomdomain.Repository
	schedules scheduledomain.Repository
	slots     slotdomain.Repository
	clock     commondomain.Clock
}

func New(
	rooms roomdomain.Repository,
	schedules scheduledomain.Repository,
	slots slotdomain.Repository,
	clock commondomain.Clock,
) *Service {
	if clock == nil {
		clock = commondomain.SystemClock{}
	}

	return &Service{
		rooms:     rooms,
		schedules: schedules,
		slots:     slots,
		clock:     clock,
	}
}

func (s *Service) EnsureGeneratedForDate(ctx context.Context, roomID string, date time.Time) error {
	return s.ensureGeneratedForDate(ctx, roomID, date, false)
}

func (s *Service) GetByID(ctx context.Context, id string) (slotdomain.Slot, error) {
	return s.slots.GetByID(ctx, id)
}

func (s *Service) ListAvailableByRoomAndDate(ctx context.Context, roomID string, date time.Time) ([]slotdomain.Slot, error) {
	if _, err := s.rooms.GetByID(ctx, roomID); err != nil {
		return nil, err
	}

	if err := s.ensureGeneratedForDate(ctx, roomID, date, true); err != nil {
		return nil, err
	}

	slots, err := s.slots.ListAvailableByRoomAndDate(ctx, roomID, date)
	if err != nil {
		return nil, err
	}

	now := s.clock.Now()
	filtered := make([]slotdomain.Slot, 0, len(slots))
	for _, entity := range slots {
		if entity.IsPast(now) {
			continue
		}

		filtered = append(filtered, entity)
	}

	return filtered, nil
}

func (s *Service) ensureGeneratedForDate(ctx context.Context, roomID string, date time.Time, roomAlreadyChecked bool) error {
	if !roomAlreadyChecked {
		if _, err := s.rooms.GetByID(ctx, roomID); err != nil {
			return err
		}
	}

	schedule, err := s.schedules.GetByRoomID(ctx, roomID)
	if err != nil {
		if errors.Is(err, scheduledomain.ErrScheduleNotFound) {
			return nil
		}

		return err
	}

	slots := slotbuilder.BuildForDate(schedule, date)
	if len(slots) == 0 {
		return nil
	}

	return s.slots.CreateMany(ctx, slots)
}

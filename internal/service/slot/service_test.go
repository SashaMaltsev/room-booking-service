package slot

import (
	"context"
	"errors"
	"testing"
	"time"

	commondomain "github.com/SashaMaltsev/room-booking-service/internal/domain/common"
	roomdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/room"
	scheduledomain "github.com/SashaMaltsev/room-booking-service/internal/domain/schedule"
	slotdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/slot"
)

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time {
	return c.now
}

type roomRepositoryStub struct {
	getByIDFn func(ctx context.Context, id string) (roomdomain.Room, error)
}

func (s roomRepositoryStub) Create(ctx context.Context, room roomdomain.Room) (roomdomain.Room, error) {
	return room, nil
}

func (s roomRepositoryStub) GetByID(ctx context.Context, id string) (roomdomain.Room, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id)
	}

	return roomdomain.Room{}, nil
}

func (s roomRepositoryStub) List(ctx context.Context) ([]roomdomain.Room, error) {
	return nil, nil
}

type scheduleRepositoryStub struct {
	getByRoomIDFn func(ctx context.Context, roomID string) (scheduledomain.Schedule, error)
}

func (s scheduleRepositoryStub) Create(ctx context.Context, schedule scheduledomain.Schedule) (scheduledomain.Schedule, error) {
	return schedule, nil
}

func (s scheduleRepositoryStub) GetByRoomID(ctx context.Context, roomID string) (scheduledomain.Schedule, error) {
	if s.getByRoomIDFn != nil {
		return s.getByRoomIDFn(ctx, roomID)
	}

	return scheduledomain.Schedule{}, nil
}

func (s scheduleRepositoryStub) ExistsForRoom(ctx context.Context, roomID string) (bool, error) {
	return false, nil
}

type slotRepositoryStub struct {
	createManyFn          func(ctx context.Context, slots []slotdomain.Slot) error
	getByIDFn             func(ctx context.Context, id string) (slotdomain.Slot, error)
	listAvailableByDateFn func(ctx context.Context, roomID string, date time.Time) ([]slotdomain.Slot, error)
}

func (s slotRepositoryStub) CreateMany(ctx context.Context, slots []slotdomain.Slot) error {
	if s.createManyFn != nil {
		return s.createManyFn(ctx, slots)
	}

	return nil
}

func (s slotRepositoryStub) GetByID(ctx context.Context, id string) (slotdomain.Slot, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id)
	}

	return slotdomain.Slot{}, nil
}

func (s slotRepositoryStub) ListByRoomAndRange(ctx context.Context, roomID string, start, end time.Time) ([]slotdomain.Slot, error) {
	return nil, nil
}

func (s slotRepositoryStub) ListAvailableByRoomAndDate(ctx context.Context, roomID string, date time.Time) ([]slotdomain.Slot, error) {
	if s.listAvailableByDateFn != nil {
		return s.listAvailableByDateFn(ctx, roomID, date)
	}

	return nil, nil
}

func TestListAvailableByRoomAndDateReturnsEmptyWhenScheduleMissing(t *testing.T) {
	service := New(
		roomRepositoryStub{
			getByIDFn: func(_ context.Context, id string) (roomdomain.Room, error) {
				return roomdomain.Room{ID: id}, nil
			},
		},
		scheduleRepositoryStub{
			getByRoomIDFn: func(_ context.Context, roomID string) (scheduledomain.Schedule, error) {
				return scheduledomain.Schedule{}, scheduledomain.ErrScheduleNotFound
			},
		},
		slotRepositoryStub{},
		commondomain.SystemClock{},
	)

	slots, err := service.ListAvailableByRoomAndDate(context.Background(), "room-1", time.Date(2026, time.March, 24, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(slots) != 0 {
		t.Fatalf("expected empty slot list, got %d items", len(slots))
	}
}

func TestListAvailableByRoomAndDateGeneratesAndFiltersPastSlots(t *testing.T) {
	startTime, _ := scheduledomain.NewTimeOfDay(9, 0)
	endTime, _ := scheduledomain.NewTimeOfDay(11, 0)

	var generated []slotdomain.Slot

	service := New(
		roomRepositoryStub{
			getByIDFn: func(_ context.Context, id string) (roomdomain.Room, error) {
				return roomdomain.Room{ID: id}, nil
			},
		},
		scheduleRepositoryStub{
			getByRoomIDFn: func(_ context.Context, roomID string) (scheduledomain.Schedule, error) {
				return scheduledomain.Schedule{
					ID:         "schedule-1",
					RoomID:     roomID,
					DaysOfWeek: []scheduledomain.Weekday{scheduledomain.Monday},
					StartTime:  startTime,
					EndTime:    endTime,
				}, nil
			},
		},
		slotRepositoryStub{
			createManyFn: func(_ context.Context, slots []slotdomain.Slot) error {
				generated = append(generated, slots...)
				return nil
			},
			listAvailableByDateFn: func(_ context.Context, roomID string, date time.Time) ([]slotdomain.Slot, error) {
				return []slotdomain.Slot{
					{
						ID:     "past",
						RoomID: roomID,
						Start:  time.Date(2026, time.March, 23, 9, 0, 0, 0, time.UTC),
						End:    time.Date(2026, time.March, 23, 9, 30, 0, 0, time.UTC),
					},
					{
						ID:     "future",
						RoomID: roomID,
						Start:  time.Date(2026, time.March, 23, 10, 0, 0, 0, time.UTC),
						End:    time.Date(2026, time.March, 23, 10, 30, 0, 0, time.UTC),
					},
				}, nil
			},
		},
		fixedClock{now: time.Date(2026, time.March, 23, 9, 45, 0, 0, time.UTC)},
	)

	slots, err := service.ListAvailableByRoomAndDate(context.Background(), "room-1", time.Date(2026, time.March, 23, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(generated) != 4 {
		t.Fatalf("expected 4 generated slots, got %d", len(generated))
	}

	if len(slots) != 1 || slots[0].ID != "future" {
		t.Fatalf("expected only future slot to remain, got %+v", slots)
	}
}

func TestEnsureGeneratedForDateReturnsRoomNotFound(t *testing.T) {
	service := New(
		roomRepositoryStub{
			getByIDFn: func(_ context.Context, id string) (roomdomain.Room, error) {
				return roomdomain.Room{}, roomdomain.ErrRoomNotFound
			},
		},
		scheduleRepositoryStub{},
		slotRepositoryStub{},
		commondomain.SystemClock{},
	)

	err := service.EnsureGeneratedForDate(context.Background(), "room-404", time.Now())
	if !errors.Is(err, roomdomain.ErrRoomNotFound) {
		t.Fatalf("expected ErrRoomNotFound, got %v", err)
	}
}

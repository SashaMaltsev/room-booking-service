package schedule

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
	createFn        func(ctx context.Context, schedule scheduledomain.Schedule) (scheduledomain.Schedule, error)
	getByRoomIDFn   func(ctx context.Context, roomID string) (scheduledomain.Schedule, error)
	existsForRoomFn func(ctx context.Context, roomID string) (bool, error)
}

func (s scheduleRepositoryStub) Create(ctx context.Context, schedule scheduledomain.Schedule) (scheduledomain.Schedule, error) {
	if s.createFn != nil {
		return s.createFn(ctx, schedule)
	}

	return schedule, nil
}

func (s scheduleRepositoryStub) GetByRoomID(ctx context.Context, roomID string) (scheduledomain.Schedule, error) {
	if s.getByRoomIDFn != nil {
		return s.getByRoomIDFn(ctx, roomID)
	}

	return scheduledomain.Schedule{}, nil
}

func (s scheduleRepositoryStub) ExistsForRoom(ctx context.Context, roomID string) (bool, error) {
	if s.existsForRoomFn != nil {
		return s.existsForRoomFn(ctx, roomID)
	}

	return false, nil
}

type slotRepositoryStub struct {
	createManyFn func(ctx context.Context, slots []slotdomain.Slot) error
}

func (s slotRepositoryStub) CreateMany(ctx context.Context, slots []slotdomain.Slot) error {
	if s.createManyFn != nil {
		return s.createManyFn(ctx, slots)
	}

	return nil
}

func (s slotRepositoryStub) GetByID(ctx context.Context, id string) (slotdomain.Slot, error) {
	return slotdomain.Slot{}, nil
}

func (s slotRepositoryStub) ListByRoomAndRange(ctx context.Context, roomID string, start, end time.Time) ([]slotdomain.Slot, error) {
	return nil, nil
}

func (s slotRepositoryStub) ListAvailableByRoomAndDate(ctx context.Context, roomID string, date time.Time) ([]slotdomain.Slot, error) {
	return nil, nil
}

func TestCreateGeneratesInitialSlotHorizon(t *testing.T) {
	startTime, _ := scheduledomain.NewTimeOfDay(9, 0)
	endTime, _ := scheduledomain.NewTimeOfDay(10, 0)

	var generated []slotdomain.Slot

	service := New(
		scheduleRepositoryStub{
			existsForRoomFn: func(_ context.Context, roomID string) (bool, error) {
				return false, nil
			},
			createFn: func(_ context.Context, schedule scheduledomain.Schedule) (scheduledomain.Schedule, error) {
				schedule.ID = "schedule-1"
				return schedule, nil
			},
		},
		roomRepositoryStub{
			getByIDFn: func(_ context.Context, id string) (roomdomain.Room, error) {
				return roomdomain.Room{ID: id}, nil
			},
		},
		slotRepositoryStub{
			createManyFn: func(_ context.Context, slots []slotdomain.Slot) error {
				generated = append(generated, slots...)
				return nil
			},
		},
		fixedClock{now: time.Date(2026, time.March, 23, 12, 0, 0, 0, time.UTC)},
		2,
	)

	created, err := service.Create(context.Background(), scheduledomain.CreateInput{
		RoomID:     "room-1",
		DaysOfWeek: []scheduledomain.Weekday{scheduledomain.Monday, scheduledomain.Tuesday},
		StartTime:  startTime,
		EndTime:    endTime,
	})
	if err != nil {
		t.Fatalf("expected create without error, got %v", err)
	}

	if created.ID != "schedule-1" {
		t.Fatalf("expected created schedule id to be set, got %q", created.ID)
	}

	if len(generated) != 4 {
		t.Fatalf("expected 4 generated slots, got %d", len(generated))
	}
}

func TestCreateReturnsScheduleExists(t *testing.T) {
	startTime, _ := scheduledomain.NewTimeOfDay(9, 0)
	endTime, _ := scheduledomain.NewTimeOfDay(10, 0)

	service := New(
		scheduleRepositoryStub{
			existsForRoomFn: func(_ context.Context, roomID string) (bool, error) {
				return true, nil
			},
		},
		roomRepositoryStub{
			getByIDFn: func(_ context.Context, id string) (roomdomain.Room, error) {
				return roomdomain.Room{ID: id}, nil
			},
		},
		slotRepositoryStub{},
		commondomain.SystemClock{},
		30,
	)

	_, err := service.Create(context.Background(), scheduledomain.CreateInput{
		RoomID:     "room-1",
		DaysOfWeek: []scheduledomain.Weekday{scheduledomain.Monday},
		StartTime:  startTime,
		EndTime:    endTime,
	})
	if !errors.Is(err, scheduledomain.ErrScheduleExists) {
		t.Fatalf("expected ErrScheduleExists, got %v", err)
	}
}

package room

import (
	"context"
	"errors"
	"testing"

	roomdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/room"
)

type roomRepositoryStub struct {
	createFn  func(ctx context.Context, room roomdomain.Room) (roomdomain.Room, error)
	getByIDFn func(ctx context.Context, id string) (roomdomain.Room, error)
	listFn    func(ctx context.Context) ([]roomdomain.Room, error)
}

func (s roomRepositoryStub) Create(ctx context.Context, room roomdomain.Room) (roomdomain.Room, error) {
	if s.createFn != nil {
		return s.createFn(ctx, room)
	}

	return room, nil
}

func (s roomRepositoryStub) GetByID(ctx context.Context, id string) (roomdomain.Room, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id)
	}

	return roomdomain.Room{}, nil
}

func (s roomRepositoryStub) List(ctx context.Context) ([]roomdomain.Room, error) {
	if s.listFn != nil {
		return s.listFn(ctx)
	}

	return nil, nil
}

func TestCreateNormalizesRoomBeforeSaving(t *testing.T) {
	var saved roomdomain.Room

	service := New(roomRepositoryStub{
		createFn: func(_ context.Context, room roomdomain.Room) (roomdomain.Room, error) {
			saved = room
			return room, nil
		},
	})

	description := "   "
	capacity := 8

	_, err := service.Create(context.Background(), roomdomain.CreateInput{
		Name:        "  Focus room  ",
		Description: &description,
		Capacity:    &capacity,
	})
	if err != nil {
		t.Fatalf("expected create without error, got %v", err)
	}

	if saved.Name != "Focus room" {
		t.Fatalf("expected trimmed room name, got %q", saved.Name)
	}

	if saved.Description != nil {
		t.Fatalf("expected empty description to become nil")
	}
}

func TestCreateReturnsValidationError(t *testing.T) {
	service := New(roomRepositoryStub{})

	capacity := 0
	_, err := service.Create(context.Background(), roomdomain.CreateInput{
		Name:     "Boardroom",
		Capacity: &capacity,
	})
	if !errors.Is(err, roomdomain.ErrCapacityMustBePositive) {
		t.Fatalf("expected ErrCapacityMustBePositive, got %v", err)
	}
}

package room

import (
	"context"

	roomdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/room"
)

var _ roomdomain.Service = (*Service)(nil)

type Service struct {
	repository roomdomain.Repository
}

func New(repository roomdomain.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(ctx context.Context, input roomdomain.CreateInput) (roomdomain.Room, error) {
	entity := roomdomain.Room{
		Name:        input.Name,
		Description: input.Description,
		Capacity:    input.Capacity,
	}
	entity.Normalize()

	if err := entity.Validate(); err != nil {
		return roomdomain.Room{}, err
	}

	return s.repository.Create(ctx, entity)
}

func (s *Service) GetByID(ctx context.Context, id string) (roomdomain.Room, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]roomdomain.Room, error) {
	return s.repository.List(ctx)
}

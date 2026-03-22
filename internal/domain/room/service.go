package room

import "context"

type CreateInput struct {
	Name        string
	Description *string
	Capacity    *int
}

type Service interface {
	Create(ctx context.Context, input CreateInput) (Room, error)
	GetByID(ctx context.Context, id string) (Room, error)
	List(ctx context.Context) ([]Room, error)
}

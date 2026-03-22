package room

import "context"

type Repository interface {
	Create(ctx context.Context, room Room) (Room, error)
	GetByID(ctx context.Context, id string) (Room, error)
	List(ctx context.Context) ([]Room, error)
}

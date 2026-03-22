package user

import (
	"context"

	"github.com/SashaMaltsev/room-booking-service/internal/domain/common"
)

type RegisterInput struct {
	Email    string
	Password string
	Role     common.Role
}

type LoginInput struct {
	Email    string
	Password string
}

type Service interface {
	Register(ctx context.Context, input RegisterInput) (User, error)
	Login(ctx context.Context, input LoginInput) (string, error)
	DummyLogin(ctx context.Context, role common.Role) (string, error)
	GetByID(ctx context.Context, id string) (User, error)
}

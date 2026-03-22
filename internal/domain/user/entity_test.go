package user

import (
	"errors"
	"testing"

	"github.com/SashaMaltsev/room-booking-service/internal/domain/common"
)

func TestUserValidateAcceptsNormalizedEmail(t *testing.T) {
	entity := User{
		Email: " USER@example.com ",
		Role:  common.RoleUser,
	}
	entity.Normalize()

	if err := entity.Validate(); err != nil {
		t.Fatalf("expected valid user, got %v", err)
	}

	if entity.Email != "user@example.com" {
		t.Fatalf("expected normalized email, got %q", entity.Email)
	}
}

func TestUserValidateRejectsInvalidEmail(t *testing.T) {
	entity := User{
		Email: "not-an-email",
		Role:  common.RoleUser,
	}

	err := entity.Validate()
	if !errors.Is(err, ErrInvalidEmail) {
		t.Fatalf("expected ErrInvalidEmail, got %v", err)
	}
}

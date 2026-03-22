package user

import (
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/SashaMaltsev/room-booking-service/internal/domain/common"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Role         common.Role
	CreatedAt    time.Time
}

func (u *User) Normalize() {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
}

func (u User) Validate() error {
	email := strings.ToLower(strings.TrimSpace(u.Email))
	if email == "" {
		return ErrEmailRequired
	}

	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email {
		return fmt.Errorf("%w: %q", ErrInvalidEmail, u.Email)
	}

	if !u.Role.IsValid() {
		return fmt.Errorf("%w: %q", common.ErrInvalidRole, u.Role)
	}

	return nil
}

func (u User) ValidateForRegistration() error {
	if err := u.Validate(); err != nil {
		return err
	}

	if strings.TrimSpace(u.PasswordHash) == "" {
		return ErrPasswordHashRequired
	}

	return nil
}

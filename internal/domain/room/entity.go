package room

import (
	"fmt"
	"strings"
	"time"
)

type Room struct {
	ID          string
	Name        string
	Description *string
	Capacity    *int
	CreatedAt   time.Time
}

func (r *Room) Normalize() {
	r.Name = strings.TrimSpace(r.Name)
	r.Description = normalizeNullableString(r.Description)
}

func (r Room) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return ErrNameRequired
	}

	if r.Capacity != nil && *r.Capacity <= 0 {
		return fmt.Errorf("%w: %d", ErrCapacityMustBePositive, *r.Capacity)
	}

	return nil
}

func normalizeNullableString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

package room

import (
	"errors"
	"testing"
)

func TestRoomNormalizeTrimsAndDropsEmptyDescription(t *testing.T) {
	description := "   "
	entity := Room{
		Name:        "  Focus room  ",
		Description: &description,
	}

	entity.Normalize()

	if entity.Name != "Focus room" {
		t.Fatalf("expected trimmed room name, got %q", entity.Name)
	}

	if entity.Description != nil {
		t.Fatalf("expected empty description to become nil")
	}
}

func TestRoomValidateRejectsNonPositiveCapacity(t *testing.T) {
	capacity := 0
	entity := Room{
		Name:     "Boardroom",
		Capacity: &capacity,
	}

	err := entity.Validate()
	if !errors.Is(err, ErrCapacityMustBePositive) {
		t.Fatalf("expected ErrCapacityMustBePositive, got %v", err)
	}
}

package slot

import (
	"errors"
	"testing"
	"time"
)

func TestSlotValidateRejectsWrongDuration(t *testing.T) {
	entity := Slot{
		RoomID: "room-1",
		Start:  time.Date(2026, time.March, 23, 9, 0, 0, 0, time.UTC),
		End:    time.Date(2026, time.March, 23, 9, 45, 0, 0, time.UTC),
	}

	err := entity.Validate()
	if !errors.Is(err, ErrInvalidSlotDuration) {
		t.Fatalf("expected ErrInvalidSlotDuration, got %v", err)
	}
}

func TestSlotOverlaps(t *testing.T) {
	first := Slot{
		RoomID: "room-1",
		Start:  time.Date(2026, time.March, 23, 9, 0, 0, 0, time.UTC),
		End:    time.Date(2026, time.March, 23, 9, 30, 0, 0, time.UTC),
	}

	second := Slot{
		RoomID: "room-1",
		Start:  time.Date(2026, time.March, 23, 9, 15, 0, 0, time.UTC),
		End:    time.Date(2026, time.March, 23, 9, 45, 0, 0, time.UTC),
	}

	if !first.Overlaps(second) {
		t.Fatalf("expected slots to overlap")
	}
}

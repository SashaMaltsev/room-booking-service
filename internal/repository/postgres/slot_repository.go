package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	slotdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/slot"
)

var _ slotdomain.Repository = (*SlotRepository)(nil)

type SlotRepository struct {
	db DBTX
}

func NewSlotRepository(db DBTX) *SlotRepository {
	return &SlotRepository{db: db}
}

func (r *SlotRepository) CreateMany(ctx context.Context, slots []slotdomain.Slot) error {
	if len(slots) == 0 {
		return nil
	}

	var builder strings.Builder
	args := make([]any, 0, len(slots)*4)

	builder.WriteString(`
		INSERT INTO slots (room_id, schedule_id, start_at, end_at)
		VALUES
	`)

	for i, slot := range slots {
		if i > 0 {
			builder.WriteString(",")
		}

		base := i * 4
		builder.WriteString(fmt.Sprintf("($%d, $%d, $%d, $%d)", base+1, base+2, base+3, base+4))

		args = append(args,
			slot.RoomID,
			slot.ScheduleID,
			slot.Start.UTC(),
			slot.End.UTC(),
		)
	}

	builder.WriteString(`
		ON CONFLICT (room_id, start_at, end_at) DO NOTHING
	`)

	_, err := r.db.ExecContext(ctx, builder.String(), args...)
	return err
}

func (r *SlotRepository) GetByID(ctx context.Context, id string) (slotdomain.Slot, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, room_id, schedule_id, start_at, end_at, created_at
		FROM slots
		WHERE id = $1
	`, id)

	entity, err := scanSlot(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return slotdomain.Slot{}, slotdomain.ErrSlotNotFound
		}

		return slotdomain.Slot{}, err
	}

	return entity, nil
}

func (r *SlotRepository) ListByRoomAndRange(ctx context.Context, roomID string, start, end time.Time) ([]slotdomain.Slot, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, room_id, schedule_id, start_at, end_at, created_at
		FROM slots
		WHERE room_id = $1
			AND start_at >= $2
			AND start_at < $3
		ORDER BY start_at ASC, id ASC
	`, roomID, start.UTC(), end.UTC())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	slots := make([]slotdomain.Slot, 0)
	for rows.Next() {
		entity, scanErr := scanSlot(rows)
		if scanErr != nil {
			return nil, scanErr
		}

		slots = append(slots, entity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return slots, nil
}

func (r *SlotRepository) ListAvailableByRoomAndDate(ctx context.Context, roomID string, date time.Time) ([]slotdomain.Slot, error) {
	start, end := utcDayBounds(date)

	rows, err := r.db.QueryContext(ctx, `
		SELECT s.id, s.room_id, s.schedule_id, s.start_at, s.end_at, s.created_at
		FROM slots s
		LEFT JOIN bookings b
			ON b.slot_id = s.id
			AND b.status = 'active'
		WHERE s.room_id = $1
			AND s.start_at >= $2
			AND s.start_at < $3
			AND b.id IS NULL
		ORDER BY s.start_at ASC, s.id ASC
	`, roomID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	slots := make([]slotdomain.Slot, 0)
	for rows.Next() {
		entity, scanErr := scanSlot(rows)
		if scanErr != nil {
			return nil, scanErr
		}

		slots = append(slots, entity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return slots, nil
}

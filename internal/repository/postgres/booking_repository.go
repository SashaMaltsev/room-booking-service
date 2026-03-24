package postgres

import (
	"context"
	"database/sql"
	"time"

	bookingdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/booking"
	commondomain "github.com/SashaMaltsev/room-booking-service/internal/domain/common"
)

var _ bookingdomain.Repository = (*BookingRepository)(nil)

type BookingRepository struct {
	db DBTX
}

func NewBookingRepository(db DBTX) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(ctx context.Context, booking bookingdomain.Booking) (bookingdomain.Booking, error) {
	booking.Normalize()

	row := r.db.QueryRowContext(ctx, `
		INSERT INTO bookings (slot_id, user_id, status, conference_link, cancelled_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, slot_id, user_id, status, conference_link, created_at, cancelled_at
	`,
		booking.SlotID,
		booking.UserID,
		string(booking.Status),
		nullableString(booking.ConferenceLink),
		nullableTime(booking.CancelledAt),
	)

	return scanBooking(row)
}

func (r *BookingRepository) GetByID(ctx context.Context, id string) (bookingdomain.Booking, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, slot_id, user_id, status, conference_link, created_at, cancelled_at
		FROM bookings
		WHERE id = $1
	`, id)

	entity, err := scanBooking(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return bookingdomain.Booking{}, bookingdomain.ErrBookingNotFound
		}

		return bookingdomain.Booking{}, err
	}

	return entity, nil
}

func (r *BookingRepository) ExistsActiveBySlotID(ctx context.Context, slotID string) (bool, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM bookings
			WHERE slot_id = $1
				AND status = 'active'
		)
	`, slotID)

	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (r *BookingRepository) ListAll(ctx context.Context, page commondomain.PageRequest) ([]bookingdomain.Booking, commondomain.Page, error) {
	normalized := page.Normalize()
	if err := normalized.Validate(); err != nil {
		return nil, commondomain.Page{}, err
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM bookings
	`).Scan(&total); err != nil {
		return nil, commondomain.Page{}, err
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, slot_id, user_id, status, conference_link, created_at, cancelled_at
		FROM bookings
		ORDER BY created_at DESC, id DESC
		LIMIT $1 OFFSET $2
	`, normalized.PageSize, normalized.Offset())
	if err != nil {
		return nil, commondomain.Page{}, err
	}
	defer rows.Close()

	bookings := make([]bookingdomain.Booking, 0)
	for rows.Next() {
		entity, scanErr := scanBooking(rows)
		if scanErr != nil {
			return nil, commondomain.Page{}, scanErr
		}

		bookings = append(bookings, entity)
	}

	if err := rows.Err(); err != nil {
		return nil, commondomain.Page{}, err
	}

	return bookings, commondomain.Page{
		Page:     normalized.Page,
		PageSize: normalized.PageSize,
		Total:    total,
	}, nil
}

func (r *BookingRepository) ListUpcomingByUser(ctx context.Context, userID string, now time.Time) ([]bookingdomain.Booking, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT b.id, b.slot_id, b.user_id, b.status, b.conference_link, b.created_at, b.cancelled_at
		FROM bookings b
		INNER JOIN slots s ON s.id = b.slot_id
		WHERE b.user_id = $1
			AND s.start_at >= $2
		ORDER BY s.start_at ASC, b.id ASC
	`, userID, now.UTC())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookings := make([]bookingdomain.Booking, 0)
	for rows.Next() {
		entity, scanErr := scanBooking(rows)
		if scanErr != nil {
			return nil, scanErr
		}

		bookings = append(bookings, entity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *BookingRepository) Update(ctx context.Context, booking bookingdomain.Booking) (bookingdomain.Booking, error) {
	booking.Normalize()

	row := r.db.QueryRowContext(ctx, `
		UPDATE bookings
		SET
			status = $2,
			conference_link = $3,
			cancelled_at = $4
		WHERE id = $1
		RETURNING id, slot_id, user_id, status, conference_link, created_at, cancelled_at
	`,
		booking.ID,
		string(booking.Status),
		nullableString(booking.ConferenceLink),
		nullableTime(booking.CancelledAt),
	)

	entity, err := scanBooking(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return bookingdomain.Booking{}, bookingdomain.ErrBookingNotFound
		}

		return bookingdomain.Booking{}, err
	}

	return entity, nil
}

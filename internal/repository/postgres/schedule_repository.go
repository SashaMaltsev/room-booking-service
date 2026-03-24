package postgres

import (
	"context"
	"database/sql"

	scheduledomain "github.com/SashaMaltsev/room-booking-service/internal/domain/schedule"
)

var _ scheduledomain.Repository = (*ScheduleRepository)(nil)

type ScheduleRepository struct {
	db DBTX
}

func NewScheduleRepository(db DBTX) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) Create(ctx context.Context, schedule scheduledomain.Schedule) (scheduledomain.Schedule, error) {
	schedule.Normalize()

	row := r.db.QueryRowContext(ctx, `
		INSERT INTO schedules (room_id, days_of_week, start_time_utc, end_time_utc)
		VALUES ($1, $2::smallint[], $3::time, $4::time)
		RETURNING
			id,
			room_id,
			days_of_week::text,
			TO_CHAR(start_time_utc, 'HH24:MI'),
			TO_CHAR(end_time_utc, 'HH24:MI'),
			created_at
	`,
		schedule.RoomID,
		encodeWeekdays(schedule.DaysOfWeek),
		schedule.StartTime.String(),
		schedule.EndTime.String(),
	)

	return scanSchedule(row)
}

func (r *ScheduleRepository) GetByRoomID(ctx context.Context, roomID string) (scheduledomain.Schedule, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT
			id,
			room_id,
			days_of_week::text,
			TO_CHAR(start_time_utc, 'HH24:MI'),
			TO_CHAR(end_time_utc, 'HH24:MI'),
			created_at
		FROM schedules
		WHERE room_id = $1
	`, roomID)

	entity, err := scanSchedule(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return scheduledomain.Schedule{}, scheduledomain.ErrScheduleNotFound
		}

		return scheduledomain.Schedule{}, err
	}

	return entity, nil
}

func (r *ScheduleRepository) ExistsForRoom(ctx context.Context, roomID string) (bool, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM schedules
			WHERE room_id = $1
		)
	`, roomID)

	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

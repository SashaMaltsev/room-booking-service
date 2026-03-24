package postgres

import (
	"database/sql"

	bookingdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/booking"
	commondomain "github.com/SashaMaltsev/room-booking-service/internal/domain/common"
	roomdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/room"
	scheduledomain "github.com/SashaMaltsev/room-booking-service/internal/domain/schedule"
	slotdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/slot"
	userdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/user"
)

func scanRoom(scanner rowScanner) (roomdomain.Room, error) {
	var (
		entity      roomdomain.Room
		description sql.NullString
		capacity    sql.NullInt64
	)

	err := scanner.Scan(
		&entity.ID,
		&entity.Name,
		&description,
		&capacity,
		&entity.CreatedAt,
	)
	if err != nil {
		return roomdomain.Room{}, err
	}

	if description.Valid {
		value := description.String
		entity.Description = &value
	}

	if capacity.Valid {
		value := int(capacity.Int64)
		entity.Capacity = &value
	}

	entity.CreatedAt = entity.CreatedAt.UTC()
	return entity, nil
}

func scanUser(scanner rowScanner) (userdomain.User, error) {
	var (
		entity       userdomain.User
		passwordHash sql.NullString
		role         string
	)

	err := scanner.Scan(
		&entity.ID,
		&entity.Email,
		&passwordHash,
		&role,
		&entity.CreatedAt,
	)
	if err != nil {
		return userdomain.User{}, err
	}

	if passwordHash.Valid {
		entity.PasswordHash = passwordHash.String
	}

	entity.Role = commondomain.Role(role)
	entity.CreatedAt = entity.CreatedAt.UTC()
	return entity, nil
}

func scanSchedule(scanner rowScanner) (scheduledomain.Schedule, error) {
	var (
		entity       scheduledomain.Schedule
		daysRaw      string
		startTimeRaw string
		endTimeRaw   string
	)

	err := scanner.Scan(
		&entity.ID,
		&entity.RoomID,
		&daysRaw,
		&startTimeRaw,
		&endTimeRaw,
		&entity.CreatedAt,
	)
	if err != nil {
		return scheduledomain.Schedule{}, err
	}

	daysOfWeek, err := decodeWeekdays(daysRaw)
	if err != nil {
		return scheduledomain.Schedule{}, err
	}

	startTime, err := scheduledomain.ParseTimeOfDay(startTimeRaw)
	if err != nil {
		return scheduledomain.Schedule{}, err
	}

	endTime, err := scheduledomain.ParseTimeOfDay(endTimeRaw)
	if err != nil {
		return scheduledomain.Schedule{}, err
	}

	entity.DaysOfWeek = daysOfWeek
	entity.StartTime = startTime
	entity.EndTime = endTime
	entity.CreatedAt = entity.CreatedAt.UTC()

	return entity, nil
}

func scanSlot(scanner rowScanner) (slotdomain.Slot, error) {
	var entity slotdomain.Slot

	err := scanner.Scan(
		&entity.ID,
		&entity.RoomID,
		&entity.ScheduleID,
		&entity.Start,
		&entity.End,
		&entity.CreatedAt,
	)
	if err != nil {
		return slotdomain.Slot{}, err
	}

	entity.Start = entity.Start.UTC()
	entity.End = entity.End.UTC()
	entity.CreatedAt = entity.CreatedAt.UTC()
	return entity, nil
}

func scanBooking(scanner rowScanner) (bookingdomain.Booking, error) {
	var (
		entity         bookingdomain.Booking
		status         string
		conferenceLink sql.NullString
		cancelledAt    sql.NullTime
	)

	err := scanner.Scan(
		&entity.ID,
		&entity.SlotID,
		&entity.UserID,
		&status,
		&conferenceLink,
		&entity.CreatedAt,
		&cancelledAt,
	)
	if err != nil {
		return bookingdomain.Booking{}, err
	}

	entity.Status = bookingdomain.Status(status)
	if conferenceLink.Valid {
		value := conferenceLink.String
		entity.ConferenceLink = &value
	}

	if cancelledAt.Valid {
		value := cancelledAt.Time.UTC()
		entity.CancelledAt = &value
	}

	entity.CreatedAt = entity.CreatedAt.UTC()
	return entity, nil
}

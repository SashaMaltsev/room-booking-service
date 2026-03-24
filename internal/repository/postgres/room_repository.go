package postgres

import (
	"context"
	"database/sql"

	roomdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/room"
)

var _ roomdomain.Repository = (*RoomRepository)(nil)

type RoomRepository struct {
	db DBTX
}

func NewRoomRepository(db DBTX) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) Create(ctx context.Context, room roomdomain.Room) (roomdomain.Room, error) {
	room.Normalize()

	row := r.db.QueryRowContext(ctx, `
		INSERT INTO rooms (name, description, capacity)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, capacity, created_at
	`,
		room.Name,
		nullableString(room.Description),
		nullableInt(room.Capacity),
	)

	return scanRoom(row)
}

func (r *RoomRepository) GetByID(ctx context.Context, id string) (roomdomain.Room, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, description, capacity, created_at
		FROM rooms
		WHERE id = $1
	`, id)

	entity, err := scanRoom(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return roomdomain.Room{}, roomdomain.ErrRoomNotFound
		}

		return roomdomain.Room{}, err
	}

	return entity, nil
}

func (r *RoomRepository) List(ctx context.Context) ([]roomdomain.Room, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, description, capacity, created_at
		FROM rooms
		ORDER BY created_at ASC, id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rooms := make([]roomdomain.Room, 0)
	for rows.Next() {
		entity, scanErr := scanRoom(rows)
		if scanErr != nil {
			return nil, scanErr
		}

		rooms = append(rooms, entity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

package postgres

import (
	"context"
	"database/sql"
	"strings"

	userdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/user"
)

var _ userdomain.Repository = (*UserRepository)(nil)

type UserRepository struct {
	db DBTX
}

func NewUserRepository(db DBTX) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user userdomain.User) (userdomain.User, error) {
	user.Normalize()

	row := r.db.QueryRowContext(ctx, `
		INSERT INTO users (email, password_hash, role)
		VALUES ($1, $2, $3)
		RETURNING id, email, password_hash, role, created_at
	`,
		user.Email,
		nullIfEmpty(user.PasswordHash),
		user.Role.String(),
	)

	return scanUser(row)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (userdomain.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash, role, created_at
		FROM users
		WHERE email = $1
	`, strings.ToLower(strings.TrimSpace(email)))

	entity, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return userdomain.User{}, userdomain.ErrUserNotFound
		}

		return userdomain.User{}, err
	}

	return entity, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (userdomain.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash, role, created_at
		FROM users
		WHERE id = $1
	`, id)

	entity, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return userdomain.User{}, userdomain.ErrUserNotFound
		}

		return userdomain.User{}, err
	}

	return entity, nil
}

func nullIfEmpty(value string) any {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}

	return trimmed
}

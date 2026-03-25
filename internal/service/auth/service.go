package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	jwtauth "github.com/SashaMaltsev/room-booking-service/internal/auth"
	commondomain "github.com/SashaMaltsev/room-booking-service/internal/domain/common"
	userdomain "github.com/SashaMaltsev/room-booking-service/internal/domain/user"
)

const (
	dummyAdminID    = "00000000-0000-0000-0000-000000000001"
	dummyAdminEmail = "dummy-admin@example.com"
	dummyUserID     = "00000000-0000-0000-0000-000000000002"
	dummyUserEmail  = "dummy-user@example.com"
)

var _ userdomain.Service = (*Service)(nil)

type Service struct {
	users  userdomain.Repository
	tokens *jwtauth.Manager
}

func New(users userdomain.Repository, tokens *jwtauth.Manager) *Service {
	return &Service{
		users:  users,
		tokens: tokens,
	}
}

func (s *Service) Register(ctx context.Context, input userdomain.RegisterInput) (userdomain.User, error) {
	entity := userdomain.User{
		Email:        input.Email,
		PasswordHash: hashPassword(input.Password),
		Role:         input.Role,
	}
	entity.Normalize()

	if err := entity.ValidateForRegistration(); err != nil {
		return userdomain.User{}, err
	}

	return s.users.Create(ctx, entity)
}

func (s *Service) Login(ctx context.Context, input userdomain.LoginInput) (string, error) {
	entity, err := s.users.GetByEmail(ctx, input.Email)
	if err != nil {
		return "", err
	}

	if entity.PasswordHash == "" || entity.PasswordHash != hashPassword(input.Password) {
		return "", userdomain.ErrInvalidCredentials
	}

	return s.tokens.Issue(entity.ID, entity.Role)
}

func (s *Service) DummyLogin(ctx context.Context, role commondomain.Role) (string, error) {
	id, email, err := dummyIdentity(role)
	if err != nil {
		return "", err
	}

	entity, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if !errors.Is(err, userdomain.ErrUserNotFound) {
			return "", err
		}

		entity, err = s.users.Create(ctx, userdomain.User{
			ID:    id,
			Email: email,
			Role:  role,
		})
		if err != nil {
			// In case of a concurrent create, try to load the same dummy user again.
			entity, err = s.users.GetByEmail(ctx, email)
			if err != nil {
				return "", err
			}
		}
	}

	return s.tokens.Issue(entity.ID, entity.Role)
}

func (s *Service) GetByID(ctx context.Context, id string) (userdomain.User, error) {
	return s.users.GetByID(ctx, id)
}

func dummyIdentity(role commondomain.Role) (id string, email string, err error) {
	switch role {
	case commondomain.RoleAdmin:
		return dummyAdminID, dummyAdminEmail, nil
	case commondomain.RoleUser:
		return dummyUserID, dummyUserEmail, nil
	default:
		return "", "", commondomain.ErrInvalidRole
	}
}

func hashPassword(password string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(password)))
	return hex.EncodeToString(sum[:])
}

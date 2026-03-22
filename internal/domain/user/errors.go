package user

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrEmailRequired        = errors.New("email is required")
	ErrInvalidEmail         = errors.New("email is invalid")
	ErrPasswordHashRequired = errors.New("password hash is required")
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrInvalidCredentials   = errors.New("invalid credentials")
)
